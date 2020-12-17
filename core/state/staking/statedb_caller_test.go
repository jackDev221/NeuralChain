package staking_test

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lvbin2012/NeuralChain/accounts/abi/bind"
	"github.com/lvbin2012/NeuralChain/accounts/abi/bind/backends"
	"github.com/lvbin2012/NeuralChain/common"
	"github.com/lvbin2012/NeuralChain/consensus/staking_contracts"
	"github.com/lvbin2012/NeuralChain/core"
	"github.com/lvbin2012/NeuralChain/crypto"
)

func TestCheckIndex(t *testing.T) {
	var (
		a, _       = common.NeutAddressStringToAddressCheck("NTkhwcpZULSbKURKqw3PYV5GEbhZFXjjBK")
		b, _       = common.NeutAddressStringToAddressCheck("NZXRfVCDbp8yttymzTg1FZ3Z4c5eJiKPDk")
		c, _       = common.NeutAddressStringToAddressCheck("NZVbXC7g42pC3tTDXUSAC5zDxGD8xpvrgJ")
		candidates = []common.Address{
			a,
			b,
		}
		epoch             = big.NewInt(40)
		startBlock        = big.NewInt(1)
		maxValidatorSize  = big.NewInt(100)
		minValidatorStake = big.NewInt(20)
		minVoteCap        = big.NewInt(10)
		adminAddr         = c
	)
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err)
	publicKey := privateKey.Public()
	addr := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	be := backends.NewSimulatedBackend(core.GenesisAlloc{
		addr: core.GenesisAccount{
			Balance: big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil),
		},
	}, gasLimit)

	authOpts := bind.NewKeyedTransactor(privateKey)
	authOpts.Nonce = big.NewInt(0)

	scAddress, tx, _, err := staking_contracts.DeployStakingContracts(authOpts, be, candidates, candidates, epoch, startBlock, maxValidatorSize, minValidatorStake, minVoteCap, adminAddr)
	require.NoError(t, err)

	be.Commit()

	receipt, err := be.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status)

	stateDB, err := be.CurrentStateDb()
	require.NoError(t, err)

	// startBlock 5
	startBlockData := stateDB.GetState(scAddress, common.BigToHash(new(big.Int).SetUint64(5)))
	assert.Equal(t, startBlockData.Big(), startBlock)

	// epoch 6
	epochData := stateDB.GetState(scAddress, common.BigToHash(new(big.Int).SetUint64(6)))
	assert.Equal(t, epochData.Big(), epoch)

	// maxValidatorSize 7
	maxValidatorSizeData := stateDB.GetState(scAddress, common.BigToHash(new(big.Int).SetUint64(7)))
	assert.Equal(t, maxValidatorSizeData.Big(), maxValidatorSize)

	// minValidatorStake 8
	minValidatorStakeData := stateDB.GetState(scAddress, common.BigToHash(new(big.Int).SetUint64(8)))
	assert.Equal(t, minValidatorStakeData.Big(), minValidatorStake)

	// minVoteCap 9
	minVoteCapData := stateDB.GetState(scAddress, common.BigToHash(new(big.Int).SetUint64(9)))
	assert.Equal(t, minVoteCapData.Big(), minVoteCap)
}
