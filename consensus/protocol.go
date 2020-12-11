// Copyright 2017 The NeuralChain Authors
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

// Package consensus implements different Evrynet consensus engines.
package consensus

import (
	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/core/types"
)

const (
	// TendermintMsg is the new message belong to neut/64.
	// it notify the protocol handler that this is a message for tendermint consensus purpose
	TendermintMsg = 0x11
)

// Broadcaster defines the interface to enqueue blocks to fetcher and find peer
type Broadcaster interface {
	// FindPeers retrives peers by addresses
	FindPeers(map[common.Address]bool) map[common.Address]Peer
	// Enqueue add a block into fetcher queue
	Enqueue(id string, block *types.Block)
}

// Peer defines the interface to communicate with peer
type Peer interface {
	// Send sends the message to this peer
	Send(msgcode uint64, data interface{}) error
	// Address return the address of a peer
	Address() common.Address
}
