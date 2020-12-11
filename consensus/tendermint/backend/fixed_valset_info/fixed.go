package fixed_valset_info

import (
	"math/big"

	"github.com/lvbin2012/NeuralChain/common"
	"github.com/lvbin2012/NeuralChain/consensus"
	"github.com/lvbin2012/NeuralChain/consensus/tendermint"
	"github.com/lvbin2012/NeuralChain/consensus/tendermint/validator"
)

type FixedValidatorSetInfo struct {
	addresses []common.Address
}

func NewFixedValidatorSetInfo(addrs []common.Address) *FixedValidatorSetInfo {
	return &FixedValidatorSetInfo{
		addresses: addrs,
	}
}

//GetValSet keep tracks of validator set in associate with blockNumber
func (mvi *FixedValidatorSetInfo) GetValSet(chainReader consensus.ChainReader, blockNumber *big.Int) (tendermint.ValidatorSet, error) {
	return validator.NewSet(mvi.addresses, tendermint.RoundRobin, blockNumber.Int64()), nil
}
