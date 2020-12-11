package tests_utils

import (
	"github.com/Evrynetlabs/evrynet-node/common"
)

// MockPeer implements consensus/protocol/Peers
type MockPeer struct {
	SendFn func(data interface{}) error
}

func (p *MockPeer) Send(msgCode uint64, data interface{}) error {
	if p.SendFn != nil {
		return p.SendFn(data)
	}
	return nil
}

func (p *MockPeer) Address() common.Address {
	return common.Address{}
}
