package tendermint

import (
	"math/big"

	"github.com/Evrynetlabs/evrynet-node/core/types"
)

// NewBlockEvent is the event sent from Backend to Core after engine.Seal() is called.
// It included the latest eligible block from tx_pool
type NewBlockEvent struct {
	Block *types.Block
}

// MessageEvent is posted for Tendermint engine communication
type MessageEvent struct {
	Payload []byte
}

// FinalCommittedEvent is posted when a proposal is committed
type FinalCommittedEvent struct {
	BlockNumber *big.Int
}

// StopCoreEvent is posted when core is stopped
type StopCoreEvent struct{}
