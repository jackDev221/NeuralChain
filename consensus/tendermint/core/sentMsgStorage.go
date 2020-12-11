package core

import (
	"io"
	"sync"

	"go.uber.org/zap"
)

// msgStorage is the struct of SOS message
type msgStorage struct {
	savedMsg []*MsgStorageData
	mu       sync.Mutex
}

// MsgStorageData contain data for message stored
type MsgStorageData struct {
	Step  RoundStepType
	Round int64
	Data  []byte
}

// NewMsgStorage returns new instance of msgStorage
func NewMsgStorage() *msgStorage {
	return &msgStorage{
		savedMsg: []*MsgStorageData{},
	}
}

// storeSentMsg stores vote/ propose to database
func (c *msgStorage) storeSentMsg(logger *zap.SugaredLogger, step RoundStepType, round int64, msgData []byte) {
	logger = logger.With("msg_step", step.String(), "msg_round", round)
	if len(msgData) == 0 {
		logger.Panicw("length of payload data should not be zero")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	logger.Debugw("saving a sent Msg")
	// Core's state is increasing step by step so the messages are inserted in correct order
	// If new msg is before last msg then logs it and returns
	if len(c.savedMsg) != 0 {
		lastData := c.savedMsg[len(c.savedMsg)-1]
		if lastData.Round > round || (lastData.Round == round && lastData.Step > step) {
			logger.Errorw("message data is before last save data, skipping")
			return
		}
	}
	c.savedMsg = append(c.savedMsg, &MsgStorageData{
		Round: round,
		Step:  step,
		Data:  msgData,
	})
}

// lookupSentMsg lockups proposal/ vote messages had stored and return index of the message
// if message were not found returns -1
func (c *msgStorage) lookup(step RoundStepType, round int64) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, element := range c.savedMsg {
		if element.Round > round {
			return i
		}
		if element.Round == round && element.Step >= step {
			return i
		}
	}
	return -1
}

func (c *msgStorage) get(index int) ([]byte, error) {
	if index >= len(c.savedMsg) || index < 0 {
		return nil, io.EOF
	}
	return c.savedMsg[index].Data, nil
}

// truncateMsgStored removes all data stored by the block's number
func (c *msgStorage) truncateMsgStored(logger *zap.SugaredLogger) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.savedMsg = []*MsgStorageData{}
	logger.Infow("truncate msgStorage done")
}
