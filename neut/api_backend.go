// Copyright 2015 The NeuralChain Authors
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

package neut

import (
	"context"
	"errors"
	"math/big"

	"github.com/lvbin2012/NeuralChain/accounts"
	"github.com/lvbin2012/NeuralChain/common"
	"github.com/lvbin2012/NeuralChain/common/math"
	"github.com/lvbin2012/NeuralChain/core"
	"github.com/lvbin2012/NeuralChain/core/bloombits"
	"github.com/lvbin2012/NeuralChain/core/rawdb"
	"github.com/lvbin2012/NeuralChain/core/state"
	"github.com/lvbin2012/NeuralChain/core/types"
	"github.com/lvbin2012/NeuralChain/core/vm"
	"github.com/lvbin2012/NeuralChain/event"
	"github.com/lvbin2012/NeuralChain/neut/downloader"
	"github.com/lvbin2012/NeuralChain/neut/gasprice"
	"github.com/lvbin2012/NeuralChain/neutdb"
	"github.com/lvbin2012/NeuralChain/params"
	"github.com/lvbin2012/NeuralChain/rpc"
)

// NeutAPIBackend implements neutapi.Backend for full nodes
type NeutAPIBackend struct {
	extRPCEnabled bool
	neut          *NeuralChain
	gpo           *gasprice.Oracle
}

// ChainConfig returns the active chain configuration.
func (b *NeutAPIBackend) ChainConfig() *params.ChainConfig {
	return b.neut.blockchain.Config()
}

func (b *NeutAPIBackend) CurrentBlock() *types.Block {
	return b.neut.blockchain.CurrentBlock()
}

func (b *NeutAPIBackend) SetHead(number uint64) {
	b.neut.protocolManager.downloader.Cancel()
	b.neut.blockchain.SetHead(number)
}

func (b *NeutAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.neut.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.neut.blockchain.CurrentBlock().Header(), nil
	}
	return b.neut.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *NeutAPIBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.neut.blockchain.GetHeaderByHash(hash), nil
}

func (b *NeutAPIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.neut.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.neut.blockchain.CurrentBlock(), nil
	}
	return b.neut.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *NeutAPIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.neut.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if err != nil {
		return nil, nil, err
	}
	if header == nil {
		return nil, nil, errors.New("header not found")
	}
	stateDb, err := b.neut.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *NeutAPIBackend) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.neut.blockchain.GetBlockByHash(hash), nil
}

func (b *NeutAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return b.neut.blockchain.GetReceiptsByHash(hash), nil
}

func (b *NeutAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	receipts := b.neut.blockchain.GetReceiptsByHash(hash)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *NeutAPIBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.neut.blockchain.GetTdByHash(blockHash)
}

func (b *NeutAPIBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.neut.BlockChain(), nil)
	return vm.NewEVM(context, state, b.neut.blockchain.Config(), *b.neut.blockchain.GetVMConfig()), vmError, nil
}

func (b *NeutAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.neut.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *NeutAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.neut.BlockChain().SubscribeChainEvent(ch)
}

func (b *NeutAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.neut.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *NeutAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.neut.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *NeutAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.neut.BlockChain().SubscribeLogsEvent(ch)
}

func (b *NeutAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.neut.txPool.AddLocal(signedTx)
}

func (b *NeutAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.neut.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *NeutAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.neut.txPool.Get(hash)
}

func (b *NeutAPIBackend) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index := rawdb.ReadTransaction(b.neut.ChainDb(), txHash)
	return tx, blockHash, blockNumber, index, nil
}

func (b *NeutAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.neut.txPool.State().GetNonce(addr), nil
}

func (b *NeutAPIBackend) Stats() (pending int, queued int) {
	return b.neut.txPool.Stats()
}

func (b *NeutAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.neut.TxPool().Content()
}

func (b *NeutAPIBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.neut.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *NeutAPIBackend) Downloader() *downloader.Downloader {
	return b.neut.Downloader()
}

func (b *NeutAPIBackend) ProtocolVersion() int {
	return b.neut.EthVersion()
}

func (b *NeutAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *NeutAPIBackend) ChainDb() neutdb.Database {
	return b.neut.ChainDb()
}

func (b *NeutAPIBackend) EventMux() *event.TypeMux {
	return b.neut.EventMux()
}

func (b *NeutAPIBackend) AccountManager() *accounts.Manager {
	return b.neut.AccountManager()
}

func (b *NeutAPIBackend) ExtRPCEnabled() bool {
	return b.extRPCEnabled
}

func (b *NeutAPIBackend) RPCGasCap() *big.Int {
	return b.neut.config.RPCGasCap
}

func (b *NeutAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.neut.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *NeutAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.neut.bloomRequests)
	}
}
