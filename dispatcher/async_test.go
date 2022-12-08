package dispatcher

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestAsync_Dispatch(t *testing.T) {
	const n = 5
	testGID := getGID()

	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		Async().Dispatch(func() {
			assert.NotEqual(t, testGID, getGID())
			wg.Done()
		})
	}
	wg.Wait()
}

func TestAsync_Dispatch_Panic(t *testing.T) {
	Async().Dispatch(func() {
		panic("should not panic")
	})
}

func TestAsync_Priority(t *testing.T) {
	require.Equal(t, AsyncPriority, Async().Priority())
}
