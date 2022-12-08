package signal

import (
	"github.com/majfault/signal/dispatcher"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignal1System(t *testing.T) {
	require := require.New(t)

	expected := []int{1, 1, 2, 3, 5}
	var out []int
	slot := func(i int) {
		out = append(out, i)
	}

	// Connect
	var sut Signal1[int]
	pSlot := sut.Connect(dispatcher.Direct(), slot)
	require.Equal(1, len(sut.slots))

	// Emit
	for _, input := range expected {
		sut.Emit(input)
	}
	require.EqualValues(expected, out)

	// Disconnect
	out = nil
	sut.Connect(dispatcher.Direct(), slot)
	require.Equal(2, len(sut.slots))
	sut.Disconnect(pSlot)
	require.Equal(1, len(sut.slots))
	for _, input := range expected {
		sut.Emit(input)
	}
	require.EqualValues(expected, out)

	// Closed
	out = nil
	sut.Close()
	require.Equal(0, len(sut.slots))
	for _, input := range expected {
		sut.Emit(input)
	}
	require.Nil(out)
}
