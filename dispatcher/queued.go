package dispatcher

import (
	"context"
	"fmt"
	"os"
	"runtime"
)

// Queued creates a new 'queued' dispatcher.
// A queue (channel) and a goroutine will be created to dequeue and call slots on first dispatch.
// The dispatcher can be closed with Close() but will be restarted on Dispatch.
type Queued struct {
	fifo        chan func()
	cancelFunc  context.CancelFunc
	Cap         int  // queue capacity: a zero or negative capacity is treated like 1 (unbuffered channel, not recommended).
	BlockOnFull bool // whether to block the caller when queue is full, or drop the slot call (default)
	LockThread  bool // whether the goroutine should always run on the same thread (eg: slots modifying OpenGL state)
}

func (q *Queued) Dispatch(slot func()) {
	if q.fifo == nil {
		if q.Cap > 1 {
			q.fifo = make(chan func(), q.Cap)
		} else {
			q.fifo = make(chan func())
		}
		var ctx context.Context
		ctx, q.cancelFunc = context.WithCancel(context.Background())
		go q.restart(ctx)
	}
	if q.Len() < cap(q.fifo) || q.BlockOnFull {
		q.fifo <- slot
	}
}

func (q *Queued) Priority() int {
	return QueuedPriority
}

// Len returns the current number of queued slots.
func (q *Queued) Len() int {
	if q.fifo != nil {
		return len(q.fifo)
	}
	return 0
}

// Close terminates the goroutine dispatcher if running (closes the channel).
func (q *Queued) Close() error {
	if q.cancelFunc != nil {
		q.cancelFunc()
		q.cancelFunc = nil
	}
	return nil
}

func (q *Queued) restart(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Recovered from slot panic:", r)
			go q.restart(ctx)
		}
	}()
	if q.LockThread {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}
	for {
		select {
		case <-ctx.Done():
			fifo := q.fifo
			q.fifo = nil
			close(fifo)
			return
		case call, ok := <-q.fifo:
			if ok {
				call()
			}
		}
	}
}
