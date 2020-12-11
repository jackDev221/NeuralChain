package staking

import (
	"math/big"

	"github.com/lvbin2012/NeuralChain/common"
	"github.com/lvbin2012/NeuralChain/consensus"
	"github.com/lvbin2012/NeuralChain/consensus/tendermint"
	"github.com/lvbin2012/NeuralChain/consensus/tendermint/utils"
	"github.com/lvbin2012/NeuralChain/consensus/tendermint/validator"
	"github.com/lvbin2012/NeuralChain/log"
)

// StakingValidator is implementation of ValidatorSetInfo
type StakingValidator struct {
	Epoch          uint64
	ProposerPolicy tendermint.ProposerPolicy
}

// NewStakingValidatorInfo returns new StakingValidator
func NewStakingValidatorInfo(epoch uint64, proposerPolicy tendermint.ProposerPolicy) *StakingValidator {
	return &StakingValidator{
		Epoch:          epoch,
		ProposerPolicy: proposerPolicy,
	}
}

// GetValSet returns the validators available in the block if it already been created
func (v *StakingValidator) GetValSet(chainReader consensus.ChainReader, number *big.Int) (tendermint.ValidatorSet, error) {
	var (
		// get the checkpoint of block-number
		blockNumber = number.Int64()
		checkPoint  = utils.GetCheckpointNumber(v.Epoch, number.Uint64())
		valSet      = validator.NewSet([]common.Address{}, v.ProposerPolicy, blockNumber)
	)

	header := chainReader.GetHeaderByNumber(checkPoint)
	if (header == nil || header.Hash() == common.Hash{}) {
		return valSet, tendermint.ErrUnknownBlock
	}

	validatorAdds, err := utils.GetValSetAddresses(header)
	if err != nil {
		log.Error("can't get the validators's address from extra-data", "number", blockNumber)
		return valSet, err
	}

	return validator.NewSet(validatorAdds, v.ProposerPolicy, blockNumber), nil
}
