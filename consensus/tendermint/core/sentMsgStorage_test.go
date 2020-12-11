package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/Evrynetlabs/evrynet-node/rlp"
)

func TestCore_StoreAndLookupSentMsg(t *testing.T) {
	zap.ReplaceGlobals(zap.NewExample())
	var (
		sentMsgStorage = NewMsgStorage()
		payload        = []byte("abc")
		payload2       = []byte("def")
		payload3       = []byte("xyz")
	)

	sentMsgStorage.storeSentMsg(zap.S(), RoundStepPropose, 1, payload)
	sentMsgStorage.storeSentMsg(zap.S(), RoundStepPrevote, 1, payload)
	sentMsgStorage.storeSentMsg(zap.S(), RoundStepPropose, 2, payload)
	sentMsgStorage.storeSentMsg(zap.S(), RoundStepPrecommit, 2, payload2)
	sentMsgStorage.storeSentMsg(zap.S(), RoundStepPrecommit, 1, payload3) // this msg doesn't follow msg order should be ignore
	sentMsgStorage.storeSentMsg(zap.S(), RoundStepPropose, 3, payload)
	sentMsgStorage.storeSentMsg(zap.S(), RoundStepPrevote, 3, payload)
	sentMsgStorage.storeSentMsg(zap.S(), RoundStepPrecommit, 3, payload)

	assert.Equal(t, 0, sentMsgStorage.lookup(RoundStepPropose, 1))
	data, err := sentMsgStorage.get(0)
	require.NoError(t, err)
	assert.Equal(t, data, payload)
	assert.Equal(t, 2, sentMsgStorage.lookup(RoundStepPrecommit, 1))
	data, err = sentMsgStorage.get(2)
	require.NoError(t, err)
	assert.Equal(t, data, payload)
	assert.Equal(t, 2, sentMsgStorage.lookup(RoundStepPropose, 2))
	assert.Equal(t, 3, sentMsgStorage.lookup(RoundStepPrevote, 2))
	data, err = sentMsgStorage.get(3)
	require.NoError(t, err)
	assert.Equal(t, data, payload2)
	assert.Equal(t, -1, sentMsgStorage.lookup(RoundStepPrevote, 4))
}

func TestCore_TruncateSentMsg(t *testing.T) {
	zap.ReplaceGlobals(zap.NewExample())
	var (
		sentMsgStorage = NewMsgStorage()
		step           = RoundStepPropose
		round          = int64(0)
		payload        []byte
	)
	msgData, _ := rlp.EncodeToBytes(payload)
	sentMsgStorage.storeSentMsg(zap.S(), step, round, msgData)

	index := sentMsgStorage.lookup(step, round)
	assert.Equal(t, 0, index)

	sentMsgStorage.truncateMsgStored(zap.S())
	index = sentMsgStorage.lookup(step, round)
	assert.Equal(t, -1, index)
}
