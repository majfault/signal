package signal

import (
	"github.com/majfault/signal/dispatcher"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignal3System(t *testing.T) {
	require := require.New(t)

	expected := []int{1, 1, 2, 3, 5}
	var out []int
	slot := func(i int, b bool, s string) {
		require.True(b)
		require.Equal("string", s)
		out = append(out, i)
	}

	// Connect
	var sut Signal3[int, bool, string]
	pSlot := sut.Connect(dispatcher.Direct(), slot)
	require.Equal(1, len(sut.slots))

	// Emit
	for _, input := range expected {
		sut.Emit(input, true, "string")
	}
	require.EqualValues(expected, out)

	// Disconnect
	out = nil
	sut.Connect(dispatcher.Direct(), slot)
	require.Equal(2, len(sut.slots))
	sut.Disconnect(pSlot)
	require.Equal(1, len(sut.slots))
	for _, input := range expected {
		sut.Emit(input, true, "string")
	}
	require.EqualValues(expected, out)

	// Closed
	out = nil
	sut.Close()
	require.Equal(0, len(sut.slots))
	for _, input := range expected {
		sut.Emit(input, true, "string")
	}
	require.Nil(out)
}
