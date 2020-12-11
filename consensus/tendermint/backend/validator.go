package backend

import (
	"math/big"

	"github.com/lvbin2012/NeuralChain/consensus"
	"github.com/lvbin2012/NeuralChain/consensus/tendermint"
)

//ValidatorSetInfo keep tracks of validator set in associate with blockNumber
type ValidatorSetInfo interface {
	GetValSet(chainReader consensus.ChainReader, blockNumber *big.Int) (tendermint.ValidatorSet, error)
}
