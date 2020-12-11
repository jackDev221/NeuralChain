package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Evrynetlabs/evrynet-node/rlp"
)

func TestCatchUpRequestMsg_DecodeRLP(t *testing.T) {
	msg := CatchUpRequestMsg{
		BlockNumber: big.NewInt(5),
		Round:       3,
		Step:        RoundStepPrecommit,
	}
	data, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)
	var decodedMsg CatchUpRequestMsg
	require.NoError(t, rlp.DecodeBytes(data, &decodedMsg))
	require.Equal(t, msg, decodedMsg)
}

func TestCatchUpReplyMsg_DecodeRLP(t *testing.T) {
	var (
		payload1 = []byte("aaaaa")
		payload2 = []byte("cc")
	)
	catchUpReplyMsg := CatchUpReplyMsg{
		Payloads:    append([][]byte{}, payload1, payload2),
		BlockNumber: big.NewInt(24),
	}
	bs, err := rlp.EncodeToBytes(&catchUpReplyMsg)
	require.NoError(t, err)
	var newMsg CatchUpReplyMsg
	require.NoError(t, rlp.DecodeBytes(bs, &newMsg))

	require.Equal(t, uint64(24), newMsg.BlockNumber.Uint64())
	require.Equal(t, payload1, newMsg.Payloads[0])
	require.Equal(t, payload2, newMsg.Payloads[1])
}
