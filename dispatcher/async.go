package dispatcher

import (
	"fmt"
	"os"
)

// Async returns the global 'asynchronous' dispatcher.
// Each slot is called in a new goroutine.
// Call ordering is not guaranteed.
func Async() Interface {
	return asyncDispatcher
}

type async struct{}

var asyncDispatcher = &async{}

func (a *async) Dispatch(slot func()) {
	go func(slot func()) {
		defer func() {
			if r := recover(); r != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Recovered from slot panic:", r)
			}
		}()
		slot()
	}(slot)
}

func (a *async) Priority() int {
	return AsyncPriority
}
