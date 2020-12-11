package tests_utils

import (
	"github.com/lvbin2012/NeuralChain/common"
	"github.com/lvbin2012/NeuralChain/core"
	"github.com/lvbin2012/NeuralChain/core/state"
	"github.com/lvbin2012/NeuralChain/core/types"
	"github.com/lvbin2012/NeuralChain/event"
)

//MockBlockChain is mock struct for block chain
type MockBlockChain struct {
	Statedb          *state.StateDB
	GasLimit         uint64
	ChainHeadFeed    *event.Feed
	MockCurrentBlock *types.Block
}

func (bc *MockBlockChain) CurrentBlock() *types.Block {
	if bc.MockCurrentBlock != nil {
		return bc.MockCurrentBlock
	}
	return types.NewBlock(&types.Header{
		GasLimit: bc.GasLimit,
	}, nil, nil, nil)
}

func (bc *MockBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *MockBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.Statedb, nil
}

func (bc *MockBlockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return bc.ChainHeadFeed.Subscribe(ch)
}
