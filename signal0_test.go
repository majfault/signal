package signal

import (
	"github.com/majfault/signal/dispatcher"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignal0System(t *testing.T) {
	require := require.New(t)
	var calls int
	slot := func() {
		calls++
	}

	// Connect
	var sut Signal0
	pSlot := sut.Connect(dispatcher.Direct(), slot)
	require.Equal(1, len(sut.slots))

	// Emit
	sut.Emit()
	sut.Emit()
	require.EqualValues(2, calls)

	// Disconnect
	calls = 0
	sut.Connect(dispatcher.Direct(), slot)
	require.Equal(2, len(sut.slots))
	sut.Disconnect(pSlot)
	require.Equal(1, len(sut.slots))
	sut.Emit()
	require.EqualValues(1, calls)

	// Closed
	calls = 0
	sut.Close()
	require.Equal(0, len(sut.slots))
	sut.Emit()
	require.Equal(0, calls)
}
