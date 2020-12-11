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

package tests

import (
	"fmt"

	"github.com/lvbin2012/NeuralChain/common"
	"github.com/lvbin2012/NeuralChain/common/hexutil"
	"github.com/lvbin2012/NeuralChain/core"
	"github.com/lvbin2012/NeuralChain/core/types"
	"github.com/lvbin2012/NeuralChain/params"
	"github.com/lvbin2012/NeuralChain/rlp"
)

// TransactionTest checks RLP decoding and sender derivation of transactions.
type TransactionTest struct {
	RLP       hexutil.Bytes `json:"rlp"`
	Vierville ttFork
}

type ttFork struct {
	Sender common.UnprefixedAddress `json:"sender"`
	Hash   common.UnprefixedHash    `json:"hash"`
}

func (tt *TransactionTest) Run(config *params.ChainConfig) error {

	validateTx := func(rlpData hexutil.Bytes, signer types.Signer) (*common.Address, *common.Hash, error) {
		tx := new(types.Transaction)
		if err := rlp.DecodeBytes(rlpData, tx); err != nil {
			return nil, nil, err
		}
		sender, err := types.Sender(signer, tx)
		if err != nil {
			return nil, nil, err
		}
		// Intrinsic gas
		requiredGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil)
		if err != nil {
			return nil, nil, err
		}
		if requiredGas > tx.Gas() {
			return nil, nil, fmt.Errorf("insufficient gas ( %d < %d )", tx.Gas(), requiredGas)
		}
		h := tx.Hash()
		return &sender, &h, nil
	}

	for _, testcase := range []struct {
		name   string
		signer types.Signer
		fork   ttFork
		isBase bool
	}{
		{"Vierville", types.NewOmahaSigner(config.ChainID), tt.Vierville, true},
	} {
		sender, txhash, err := validateTx(tt.RLP, testcase.signer)

		if testcase.fork.Sender == (common.UnprefixedAddress{}) {
			if err == nil {
				return fmt.Errorf("Expected error, got none (address %v)", sender.String())
			}
			continue
		}
		// Should resolve the right address
		if err != nil {
			return fmt.Errorf("Got error, expected none: %v", err)
		}
		if sender == nil {
			return fmt.Errorf("sender was nil, should be %x", common.Address(testcase.fork.Sender))
		}
		if *sender != common.Address(testcase.fork.Sender) {
			return fmt.Errorf("Sender mismatch: got %x, want %x", sender, testcase.fork.Sender)
		}
		if txhash == nil {
			return fmt.Errorf("txhash was nil, should be %x", common.Hash(testcase.fork.Hash))
		}
		if *txhash != common.Hash(testcase.fork.Hash) {
			return fmt.Errorf("Hash mismatch: got %x, want %x", *txhash, testcase.fork.Hash)
		}
	}
	return nil
}
