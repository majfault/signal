package dispatcher

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestDirect_Dispatch(t *testing.T) {
	const n = 5
	testGID := getGID()

	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		Direct().Dispatch(func() {
			require.Equal(t, testGID, getGID())
			wg.Done()
		})
	}
	wg.Wait()
}

func TestDirect_Dispatch_Panic(t *testing.T) {
	Direct().Dispatch(func() {
		panic("should not panic")
	})
}

func TestDirect_Priority(t *testing.T) {
	require.Equal(t, DirectPriority, Direct().Priority())
}
