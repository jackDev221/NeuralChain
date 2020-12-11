package core

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeoutInfo(t *testing.T) {
	ti1 := timeoutInfo{
		BlockNumber: big.NewInt(209),
		Round:       1,
		Step:        RoundStepPropose,
		Duration:    time.Second,
		Retry:       0,
	}

	ti2 := timeoutInfo{
		BlockNumber: big.NewInt(209),
		Round:       0,
		Step:        RoundStepPrecommitWait,
		Duration:    time.Second,
		Retry:       0,
	}

	ti3 := timeoutInfo{
		BlockNumber: big.NewInt(209),
		Round:       1,
		Step:        RoundStepPrecommit,
		Duration:    time.Second,
		Retry:       0,
	}

	ti4 := timeoutInfo{
		BlockNumber: big.NewInt(208),
		Round:       1,
		Step:        RoundStepPrecommit,
		Duration:    time.Second,
		Retry:       0,
	}

	require.Equal(t, false, ti1.earlierOrEqual(ti2))
	require.Equal(t, true, ti1.earlierOrEqual(ti3))
	require.Equal(t, false, ti1.earlierOrEqual(ti4))
}

func TestRunningTimeoutTicker(t *testing.T) {
	ticker := NewTimeoutTicker()
	require.NoError(t, ticker.Start())

	ticker.ScheduleTimeout(timeoutInfo{
		Duration:    time.Millisecond * 10,
		BlockNumber: big.NewInt(1),
		Round:       0,
		Step:        RoundStepPrevote,
	})
	time.Sleep(time.Millisecond * 20)

	require.NoError(t, ticker.Stop())

	//This timeoutInfo should not be sent to tokChan
	ticker.ScheduleTimeout(timeoutInfo{
		Duration:    time.Millisecond * 10,
		BlockNumber: big.NewInt(2),
		Round:       0,
		Step:        RoundStepPrevoteWait,
	})
	time.Sleep(time.Millisecond * 20)

	ti1, ok := <-ticker.Chan()
	assert.NotNil(t, ti1, "Only the first timeoutInfo can get from tokChan")
	assert.True(t, ok)
	assert.Equal(t, big.NewInt(1), ti1.BlockNumber, "The blocknumber of the first timeoutInfo must be 1")

	_, ok = <-ticker.Chan()
	assert.False(t, ok, "The second timeouInfo won't be sent to tokChan")
}

func TestRunningTimeoutTickerAbortSending(t *testing.T) {
	ticker := NewTimeoutTicker()
	require.NoError(t, ticker.Start())

	for i := 0; i < 11; i++ {
		ticker.ScheduleTimeout(timeoutInfo{
			Duration:    time.Millisecond * 1,
			BlockNumber: big.NewInt(1),
			Round:       int64(i),
			Step:        RoundStepPrevote,
		})
		time.Sleep(time.Millisecond * 2)
	}
	require.NoError(t, ticker.Stop())
	time.Sleep(time.Millisecond * 20)
}
