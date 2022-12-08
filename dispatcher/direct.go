package dispatcher

import (
	"fmt"
	"os"
)

// Direct returns the global 'direct' dispatcher.
// Each slot is called in the same context as the caller (ie like a standard function call), in the order they were connected to the signal.
// Calls are blocking for the caller.
func Direct() Interface {
	return directDispatcher
}

type direct struct{}

var directDispatcher = &direct{}

func (d *direct) Dispatch(slot func()) {
	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Recovered from slot panic:", r)
		}
	}()
	slot()
}

func (d *direct) Priority() int {
	return DirectPriority
}
