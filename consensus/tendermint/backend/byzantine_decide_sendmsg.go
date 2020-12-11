package backend

import (
	"math/rand"

	"github.com/lvbin2012/NeuralChain/consensus/tendermint"
	"github.com/lvbin2012/NeuralChain/log"
)

// checkAndSendMsg decided to send the message or not
func (sb *Backend) checkAndSendMsg(payload []byte) error {
	var decidedSendMsg = true
	if sb.config.FaultyMode == tendermint.RandomlyStopSendingMsg.Uint64() {
		// randomly stop sending message.
		switch rand.Intn(2) {
		case 0: // stop sending message
			decidedSendMsg = false
			log.Warn("Byzantine mode: stop sending message.")
		case 1: // sending message
			decidedSendMsg = true
			log.Warn("Byzantine mode: sending message.")
		}
	}

	if decidedSendMsg {
		return sb.EventMux().Post(tendermint.MessageEvent{
			Payload: payload,
		})
	}
	return nil
}
