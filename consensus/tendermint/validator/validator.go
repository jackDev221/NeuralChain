package validator

import (
	"reflect"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint"
)

// New will create new validator
func New(addr common.Address) tendermint.Validator {
	return &defaultValidator{
		address: addr,
	}
}

// NewSet will create new validator set by address list & policy
func NewSet(addrs []common.Address, policy tendermint.ProposerPolicy, height int64) tendermint.ValidatorSet {
	return newDefaultSet(addrs, policy, height)
}

// IsProposer will be checking whether the validator with given address is a proposer
func (valSet *defaultSet) IsProposer(address common.Address) bool {
	_, val := valSet.GetByAddress(address)
	return reflect.DeepEqual(valSet.GetProposer(), val)
}

// ExtractValidators will extract extra data to address list
func ExtractValidators(extraData []byte) []common.Address {
	// get the validator addresses
	addrs := make([]common.Address, len(extraData)/common.AddressLength)
	for i := 0; i < len(addrs); i++ {
		copy(addrs[i][:], extraData[i*common.AddressLength:])
	}

	return addrs
}
