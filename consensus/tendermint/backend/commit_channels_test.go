package backend

import (
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Evrynetlabs/evrynet-node/core/types"
)

func TestCommitChannels(t *testing.T) {
	var (
		blk1 int64 = 14
		blk2 int64 = 15
		ch3  <-chan *types.Block
	)
	commitChs := newCommitChannels()
	ch1 := commitChs.createCommitChannelAndCloseIfExist(strconv.FormatInt(blk1, 10))
	ch2 := commitChs.createCommitChannelAndCloseIfExist(strconv.FormatInt(blk2, 10))
	go func() {
		ch3 = commitChs.createCommitChannelAndCloseIfExist(strconv.FormatInt(blk2, 10))
	}()

	select {
	case <-ch1:
		t.Fatalf("channel should not be closed")
	case _, closed := <-ch2:
		require.False(t, closed)
	case <-time.After(time.Millisecond):
		t.Errorf("channel should be closed when new channel is created")
	}

	commitChs.sendBlock(types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(blk2),
	}))

	select {
	case <-ch1:
		t.Fatalf("channel should not be closed")
	case _, closed := <-ch3:
		require.True(t, closed)
	case <-time.After(time.Millisecond):
		t.Errorf("channel should be closed when new channel is created")
	}

}
