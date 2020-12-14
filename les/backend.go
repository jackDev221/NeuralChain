// Copyright 2016 The NeuralChain Authors
// This file is part of the NeuralChain library .
//
// The NeuralChain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The NeuralChain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the NeuralChain library . If not, see <http://www.gnu.org/licenses/>.

// Package les implements the Light NeuralChain Subprotocol.
package les

import (
	"fmt"
	"sync"
	"time"

	"github.com/lvbin2012/NeuralChain/accounts"
	"github.com/lvbin2012/NeuralChain/common"
	"github.com/lvbin2012/NeuralChain/common/hexutil"
	"github.com/lvbin2012/NeuralChain/common/mclock"
	"github.com/lvbin2012/NeuralChain/consensus"
	"github.com/lvbin2012/NeuralChain/core"
	"github.com/lvbin2012/NeuralChain/core/bloombits"
	"github.com/lvbin2012/NeuralChain/core/rawdb"
	"github.com/lvbin2012/NeuralChain/core/types"
	"github.com/lvbin2012/NeuralChain/event"
	"github.com/lvbin2012/NeuralChain/internal/neutapi"
	"github.com/lvbin2012/NeuralChain/light"
	"github.com/lvbin2012/NeuralChain/log"
	"github.com/lvbin2012/NeuralChain/neut"
	"github.com/lvbin2012/NeuralChain/neut/downloader"
	"github.com/lvbin2012/NeuralChain/neut/filters"
	"github.com/lvbin2012/NeuralChain/neut/gasprice"
	"github.com/lvbin2012/NeuralChain/node"
	"github.com/lvbin2012/NeuralChain/p2p"
	"github.com/lvbin2012/NeuralChain/p2p/discv5"
	"github.com/lvbin2012/NeuralChain/params"
	"github.com/lvbin2012/NeuralChain/rpc"
)

type LightNeuralChain struct {
	lesCommons

	odr         *LesOdr
	relay       *LesTxRelay
	chainConfig *params.ChainConfig
	// Channel for shutting down the service
	shutdownChan chan bool

	// Handlers
	peers      *peerSet
	txPool     *light.TxPool
	blockchain *light.LightChain
	serverPool *serverPool
	reqDist    *requestDistributor
	retriever  *retrieveManager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer

	ApiBackend *LesApiBackend

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	networkId     uint64
	netRPCService *neutapi.PublicNetAPI

	wg sync.WaitGroup
}

func New(ctx *node.ServiceContext, config *neut.Config) (*LightNeuralChain, error) {
	chainDb, err := ctx.OpenDatabase("lightchaindata", config.DatabaseCache, config.DatabaseHandles, "neut/db/chaindata/")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlockWithOverride(chainDb, config.Genesis)
	if _, isCompat := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !isCompat {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	peers := newPeerSet()
	quitSync := make(chan struct{})

	leth := &LightNeuralChain{
		lesCommons: lesCommons{
			chainDb: chainDb,
			config:  config,
			iConfig: light.DefaultClientIndexerConfig,
		},
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		peers:          peers,
		reqDist:        newRequestDistributor(peers, quitSync, &mclock.System{}),
		accountManager: ctx.AccountManager,
		engine:         neut.CreateConsensusEngine(ctx, chainConfig, config, nil, false, chainDb),
		shutdownChan:   make(chan bool),
		networkId:      config.NetworkId,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   neut.NewBloomIndexer(chainDb, params.BloomBitsBlocksClient, params.HelperTrieConfirmations),
	}

	var trustedNodes []string
	if leth.config.ULC != nil {
		trustedNodes = leth.config.ULC.TrustedServers
	}
	leth.serverPool = newServerPool(chainDb, quitSync, &leth.wg, trustedNodes)
	leth.retriever = newRetrieveManager(peers, leth.reqDist, leth.serverPool)
	leth.relay = NewLesTxRelay(peers, leth.retriever)

	leth.odr = NewLesOdr(chainDb, light.DefaultClientIndexerConfig, leth.retriever)
	leth.chtIndexer = light.NewChtIndexer(chainDb, leth.odr, params.CHTFrequency, params.HelperTrieConfirmations)
	leth.bloomTrieIndexer = light.NewBloomTrieIndexer(chainDb, leth.odr, params.BloomBitsBlocksClient, params.BloomTrieFrequency)
	leth.odr.SetIndexers(leth.chtIndexer, leth.bloomTrieIndexer, leth.bloomIndexer)

	// Note: NewLightChain adds the trusted checkpoint so it needs an ODR with
	// indexers already set but not started yet
	if leth.blockchain, err = light.NewLightChain(leth.odr, leth.chainConfig, leth.engine); err != nil {
		return nil, err
	}
	// Note: AddChildIndexer starts the update process for the child
	leth.bloomIndexer.AddChildIndexer(leth.bloomTrieIndexer)
	leth.chtIndexer.Start(leth.blockchain)
	leth.bloomIndexer.Start(leth.blockchain)

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		leth.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}

	leth.txPool = light.NewTxPool(leth.chainConfig, leth.blockchain, leth.relay)

	if leth.protocolManager, err = NewProtocolManager(
		leth.chainConfig,
		light.DefaultClientIndexerConfig,
		true,
		config.NetworkId,
		leth.eventMux,
		leth.engine,
		leth.peers,
		leth.blockchain,
		nil,
		chainDb,
		leth.odr,
		leth.relay,
		leth.serverPool,
		quitSync,
		&leth.wg,
		config.ULC,
		nil); err != nil {
		return nil, err
	}

	if leth.protocolManager.isULCEnabled() {
		log.Warn("Ultra light client is enabled", "trustedNodes", len(leth.protocolManager.ulc.trustedKeys), "minTrustedFraction", leth.protocolManager.ulc.minTrustedFraction)
		leth.blockchain.DisableCheckFreq()
	}
	leth.ApiBackend = &LesApiBackend{ctx.ExtRPCEnabled(), leth, nil}

	gpoParams := config.GPO
	leth.ApiBackend.gpo = gasprice.NewOracle(leth.ApiBackend, gpoParams)
	return leth, nil
}

func lesTopic(genesisHash common.Hash, protocolVersion uint) discv5.Topic {
	var name string
	switch protocolVersion {
	case lpv2:
		name = "LES2"
	default:
		panic(nil)
	}
	return discv5.Topic(name + "@" + common.Bytes2Hex(genesisHash.Bytes()[0:8]))
}

type LightDummyAPI struct{}

// Etherbase is the address that mining rewards will be send to
func (s *LightDummyAPI) Etherbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("mining is not supported in light mode")
}

// Coinbase is the address that mining rewards will be send to (alias for Etherbase)
func (s *LightDummyAPI) Coinbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("mining is not supported in light mode")
}

// Hashrate returns the POW hashrate
func (s *LightDummyAPI) Hashrate() hexutil.Uint {
	return 0
}

// Mining returns an indication if this node is currently mining.
func (s *LightDummyAPI) Mining() bool {
	return false
}

// APIs returns the collection of RPC services the neuralChain package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *LightNeuralChain) APIs() []rpc.API {
	return append(neutapi.GetAPIs(s.ApiBackend), []rpc.API{
		{
			Namespace: "neut",
			Version:   "1.0",
			Service:   &LightDummyAPI{},
			Public:    true,
		}, {
			Namespace: "neut",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "neut",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, true),
			Public:    true,
		}, {
			Namespace: "neut",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *LightNeuralChain) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *LightNeuralChain) BlockChain() *light.LightChain      { return s.blockchain }
func (s *LightNeuralChain) TxPool() *light.TxPool              { return s.txPool }
func (s *LightNeuralChain) Engine() consensus.Engine           { return s.engine }
func (s *LightNeuralChain) LesVersion() int                    { return int(ClientProtocolVersions[0]) }
func (s *LightNeuralChain) Downloader() *downloader.Downloader { return s.protocolManager.downloader }
func (s *LightNeuralChain) EventMux() *event.TypeMux           { return s.eventMux }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *LightNeuralChain) Protocols() []p2p.Protocol {
	return s.makeProtocols(ClientProtocolVersions)
}

// Start implements node.Service, starting all internal goroutines needed by the
// NeuralChain protocol implementation.
func (s *LightNeuralChain) Start(srvr *p2p.Server) error {
	log.Warn("Light client mode is an experimental feature")
	s.startBloomHandlers(params.BloomBitsBlocksClient)
	s.netRPCService = neutapi.NewPublicNetAPI(srvr, s.networkId)
	// clients are searching for the first advertised protocol in the list
	protocolVersion := AdvertiseProtocolVersions[0]
	s.serverPool.start(srvr, lesTopic(s.blockchain.Genesis().Hash(), protocolVersion))
	s.protocolManager.Start(s.config.LightPeers)
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// NeuralChain protocol.
func (s *LightNeuralChain) Stop() error {
	s.odr.Stop()
	s.relay.Stop()
	s.bloomIndexer.Close()
	s.chtIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	s.txPool.Stop()
	s.engine.Close()

	s.eventMux.Stop()

	time.Sleep(time.Millisecond * 200)
	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
