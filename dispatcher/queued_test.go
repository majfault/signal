package dispatcher

import (
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueued_Dispatch_Len_Close(t *testing.T) {
	sut := Queued{
		Cap:        8,
		LockThread: true,
	}
	require.Nil(t, sut.fifo)
	require.Equal(t, 0, sut.Len())

	fence := make(chan struct{}, 5)
	testGID := getGID()
	dispatcherGID := uint64(0)
	slowFunc := func() {
		fence <- struct{}{}
		gid := getGID()
		if atomic.CompareAndSwapUint64(&dispatcherGID, 0, gid) {
			// get initial dispatcher GID and ensure different from caller's (test goroutine)
			require.NotEqual(t, testGID, gid)
		} else {
			// ensure always in same dispatcher goroutine
			require.Equal(t, atomic.LoadUint64(&dispatcherGID), gid)
		}
		<-time.After(200 * time.Millisecond)
		fence <- struct{}{}
	}

	sut.Dispatch(slowFunc)
	sut.Dispatch(slowFunc)
	sut.Dispatch(slowFunc)

	<-fence
	require.Equal(t, 8, cap(sut.fifo))
	require.Equal(t, 2, sut.Len())

	_ = sut.Close()
	<-fence
	<-time.After(500 * time.Millisecond)
	require.Nil(t, sut.cancelFunc)
	require.Nil(t, sut.fifo)
	require.Equal(t, 0, sut.Len())
}

func TestQueued_Dispatch_Panic(t *testing.T) {
	sut := Queued{
		Cap: 8,
	}

	callCount := uint32(0)
	wg := sync.WaitGroup{}

	panicCall := func() {
		panic("should not fail")
	}
	normalCall := func() {
		atomic.AddUint32(&callCount, 1)
		wg.Done()
	}

	wg.Add(2)
	sut.Dispatch(normalCall)
	sut.Dispatch(panicCall)
	sut.Dispatch(normalCall)
	wg.Wait()

	require.Equal(t, uint32(2), callCount)
}

func TestQueued_Priority(t *testing.T) {
	sut := Queued{}
	require.Equal(t, QueuedPriority, sut.Priority())
}
