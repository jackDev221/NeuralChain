package core

import (
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/tests_utils"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/rlp"
)

func Test_CatchUpRequest(t *testing.T) {
	zap.ReplaceGlobals(zap.NewExample())
	var (
		nodePrivateKey  = tests_utils.MakeNodeKey()
		nodeAddr        = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		nodePrivateKey2 = tests_utils.MakeNodeKey()
		nodeAddr2       = crypto.PubkeyToAddress(nodePrivateKey2.PublicKey)
		validators      = []common.Address{
			nodeAddr2,
			nodeAddr,
		}
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
		tendermintCfg = tests_utils.DefaultTestConfig
	)
	//create New test backend and newMockChain
	be, _ := tests_utils.MustCreateAndStartNewBackend(t, nodePrivateKey, genesisHeader, validators)
	// subscribe to output msg
	mockBe, ok := be.(*tests_utils.MockBackend)
	require.True(t, ok)
	sentMsgSub := mockBe.SendEventMux.Subscribe(tests_utils.SentMsgEvent{})
	defer sentMsgSub.Unsubscribe()

	core := newTestCore(be, tendermintCfg)
	require.NoError(t, core.Start())
	state := core.CurrentState()
	assert.Equal(t, RoundStepNewHeight, state.Step())
	state.SetBlock(nil)

	assertToVal2 := func(address common.Address) {
		require.Equal(t, nodeAddr2, address)
	}
	assertNextMsg(t, sentMsgSub, msgPrevote, tendermintCfg.TimeoutPropose+(1*time.Second), assertToVal2, nil)
	zap.S().Debugw("assert sending catch up when being stuck to prevote")
	assertNextMsg(t, sentMsgSub, msgCatchUpRequest, tendermintCfg.TimeoutPropose*2+(1*time.Second), assertToVal2, func(Msg []byte) {
		var catchUpReq CatchUpRequestMsg
		require.NoError(t, rlp.DecodeBytes(Msg, &catchUpReq))
		assert.Equal(t, catchUpReq.Round, int64(0))
		assert.Equal(t, catchUpReq.BlockNumber, big.NewInt(1))
		assert.Equal(t, catchUpReq.Step, RoundStepPrevote)
	})

	require.NoError(t, be.EventMux().Post(tendermint.MessageEvent{
		Payload: createCatchUpRequest(t, nodePrivateKey2, big.NewInt(1), 0, RoundStepPrevote),
	}))
	zap.S().Debugw("assert handle request")
	assertNextMsg(t, sentMsgSub, msgCatchUpReply, time.Second, assertToVal2, func(Msg []byte) {
		var reply CatchUpReplyMsg
		require.NoError(t, rlp.DecodeBytes(Msg, &reply))
		assert.Equal(t, len(reply.Payloads), 1)
	})

	prevotePayload := createPrevote(t, nodePrivateKey2, big.NewInt(1), 0)
	catchUpReplyPayload := createCatchUpReply(t, nodePrivateKey2, big.NewInt(1), [][]byte{prevotePayload})
	require.NoError(t, be.EventMux().Post(tendermint.MessageEvent{
		Payload: catchUpReplyPayload,
	}))

	zap.S().Debugw("assert  handle reply")
	assertNextMsg(t, sentMsgSub, msgPrecommit, 1*time.Second, assertToVal2, nil)
}

func createCatchUpRequest(t *testing.T, privateKey *ecdsa.PrivateKey, BlockNumber *big.Int, Round int64, Step RoundStepType) []byte {
	addr := crypto.PubkeyToAddress(privateKey.PublicKey)
	msgCatchUp := CatchUpRequestMsg{
		BlockNumber: BlockNumber,
		Round:       Round,
		Step:        Step,
	}
	bs, err := rlp.EncodeToBytes(&msgCatchUp)
	require.NoError(t, err)
	msg := message{
		Address: addr,
		Msg:     bs,
		Code:    msgCatchUpRequest,
	}
	sign(t, &msg, privateKey)
	bs, err = rlp.EncodeToBytes(&msg)
	require.NoError(t, err)
	return bs
}

func createPrevote(t *testing.T, privateKey *ecdsa.PrivateKey, BlockNumber *big.Int, Round int64) []byte {
	var (
		bs        []byte
		err       error
		msg       message
		addr      = crypto.PubkeyToAddress(privateKey.PublicKey)
		emptyHash = common.Hash{}
	)
	vote := Vote{
		Round:       Round,
		BlockNumber: BlockNumber,
		BlockHash:   &emptyHash,
		Seal:        []byte{},
	}
	bs, err = rlp.EncodeToBytes(&vote)
	require.NoError(t, err)
	msg = message{
		Address: addr,
		Msg:     bs,
		Code:    msgPrevote,
	}
	sign(t, &msg, privateKey)
	prevotePayload, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)
	return prevotePayload
}

func createCatchUpReply(t *testing.T, privateKey *ecdsa.PrivateKey, BlockNumber *big.Int, payloads [][]byte) []byte {
	var (
		bs   []byte
		err  error
		msg  message
		addr = crypto.PubkeyToAddress(privateKey.PublicKey)
	)
	msgCatchUp := CatchUpReplyMsg{
		BlockNumber: BlockNumber,
		Payloads:    payloads,
	}
	bs, err = rlp.EncodeToBytes(&msgCatchUp)
	require.NoError(t, err)
	msg = message{
		Address: addr,
		Msg:     bs,
		Code:    msgCatchUpReply,
	}
	sign(t, &msg, privateKey)
	replyPayload, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)
	return replyPayload
}
