package core

import (
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint"
	"github.com/Evrynetlabs/evrynet-node/rlp"
)

func (c *core) enterCatchup(tiBlock *big.Int, tiRound int64, tiStep RoundStepType, tiRetry uint64) {
	var (
		state        = c.currentState
		sRound       = state.Round()
		sBlockNumber = state.BlockNumber()
		sStep        = state.Step()
		logger       = c.getLogger().With("ti_block", tiBlock, "ti_round", tiRound, "ti_step", tiStep)
	)

	if sBlockNumber.Cmp(tiBlock) != 0 || tiRound != sRound || sStep != tiStep {
		logger.Debugw("catchUp ignore: only sending catch-up msg if node is stuck at prevote or precommit")
		return
	}
	// schedule next catch up, this will avoid response catch-up error
	var nextCatchUpDuration time.Duration
	switch tiStep {
	case RoundStepPrevote:
		nextCatchUpDuration = c.config.PrevoteCatchupTimeout(sRound)
	case RoundStepPrecommit:
		nextCatchUpDuration = c.config.PrecommitCatchupTimeout(sRound)
	default:
		logger.Errorw("get unexpected timeout step")
		return
	}
	c.timeout.ScheduleTimeout(timeoutInfo{
		Duration:    nextCatchUpDuration,
		BlockNumber: new(big.Int).Set(sBlockNumber),
		Round:       sRound,
		Step:        tiStep,
		Retry:       tiRetry + 1,
	})
	//send catch up
	c.sendCatchUpRequest(logger, tiBlock, tiRound, tiStep)
}

func (c *core) sendCatchUpRequest(logger *zap.SugaredLogger, tiBlock *big.Int, tiRound int64, tiStep RoundStepType) {
	var (
		state = c.currentState
		addr  = c.backend.Address()
	)
	//send catch up
	msg := &CatchUpRequestMsg{
		Round:       tiRound,
		BlockNumber: new(big.Int).Set(tiBlock),
		Step:        tiStep,
	}
	msgData, err := rlp.EncodeToBytes(msg)
	if err != nil {
		logger.Errorw("Failed to encode CatchUpRequestMsg to bytes", "err", err)
		return
	}
	payload, err := c.FinalizeMsg(&message{
		Code: msgCatchUpRequest,
		Msg:  msgData,
	})
	if err != nil {
		logger.Errorw("Failed to finalize CatchUpRequestMsg to bytes", "err", err)
		return
	}

	var (
		msgSet *messageSet
		ok     bool
	)
	switch tiStep {
	case RoundStepPrevote:
		msgSet, ok = state.GetPrevotesByRound(tiRound)
	case RoundStepPrecommit:
		msgSet, ok = state.GetPrecommitsByRound(tiRound)
	default:
		logger.Errorw("get unexpected timeout step")
		return
	}
	var missing map[common.Address]bool
	if !ok { // not found any vote
		missing = make(map[common.Address]bool)
		for _, val := range c.valSet.List() {
			missing[val.Address()] = true
		}
	} else {
		missing = msgSet.MissingVotes()
	}
	// if missing votes from itself (core.Stop() before handle msg from eventMux) then take from core 's storage
	if _, ok := missing[addr]; ok {
		go func() {
			index := c.sentMsgStorage.lookup(tiStep, tiRound)
			missingPayload, err := c.sentMsgStorage.get(index)
			if err != nil {
				logger.Warnw("Failed to found self msg", err, "err")
				return
			}
			if err := c.backend.EventMux().Post(tendermint.MessageEvent{
				Payload: missingPayload,
			}); err != nil {
				logger.Errorw("Failed to re-post msg from core 's storage to eventMux")
			}
		}()
	}
	delete(missing, addr)

	if err := c.backend.Multicast(missing, payload); err != nil {
		logger.Debugw("Failed to multicast msg", "err", err.Error())
	}
}
