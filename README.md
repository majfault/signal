# Signal

[![Go Reference](https://pkg.go.dev/badge/github.com/majfault/signal)](https://pkg.go.dev/github.com/majfault/signal)
[![GitHub license](https://img.shields.io/github/license/majfault/signal)](https://github.com/majfault/signal/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/majfault/signal)](https://goreportcard.com/report/github.com/majfault/signal)

Signal/Slot event-driven paradigm using generics, for Go 1.18 and beyond...

Heavily inspired by Qt's event system and `sig4j` Java implementation.

Any number of *listeners* can `Connect()` a callback (called `Slot`) using a *dispatcher* on a `Signal`.
Every time `Emit()` is called on a `Signal`, every connected `Slot` will be called with their associated *dispatcher*.

Features include:
- Type safe event handling through the use of generics
- Many-to-many event dispatching: anyone and anything may `Emit` on and `Connect` to the same `Signal`
- *Callee*-chosen call context through *dispatchers*, as opposed to most event systems: each `Slot` is invoked using the dispatcher it was connected with.
- Several slot call dispatchers to choose from: `Direct()`, `Async()`, `Queued{}` 

## Installing

Run:
> go get -u github.com/majfault/signal

## Dispatchers

|                       | Callee/Slot call context                                                                                         | Properties                                           |
|:----------------------|:-----------------------------------------------------------------------------------------------------------------|:-----------------------------------------------------|
| `dispatcher.Direct()` | Same as caller                                                                                                   | ❌ Non-blocking for caller<br>✅ Event order garanteed |
| `dispatcher.Async()`  | New goroutine                                                                                                    | ✅ Non-blocking for caller<br>❌ Event order garanteed |
| `dispatcher.Queued{}` | Queue/FIFO with dedicated goroutine<br>Goroutine is only started on first dispatch<br>Dispatcher can be `Close()`| ✅ Non-blocking for caller<br>✅ Event order garanteed |

## Examples

### Declaring signals and emitting

```go
package main

import (
    "github.com/majfault/signal"
	"time"
)

func main()  {
	// without parameters
    sig0 := signal.Signal0{}
	
	// with 1 parameter
    sig1 := signal.Signal1[int]{}
	
	// with 2 parameters
    sig2 := signal.Signal2[string, []byte]{}
	
	// with 3 parameters
    sig3 := signal.Signal3[[]int, map[int]string, bool]{}
	
	// with 4 parameters
    sig4 := signal.Signal4[bool, int, time.Time, string]{}
	
	// emitting
    sig0.Emit()
    sig2.Emit("string", []byte("byte slice"))
    sig4.Emit(true, 42, time.Now(), "emit")
}
```

### Connecting slots

```go
package main

import (
    "github.com/majfault/signal"
    "github.com/majfault/signal/dispatcher"
    "time"    
    "fmt"
)

var keyPressed = signal.Signal1[int]{} 

func main()  {
    keyPressed.Connect(dispatcher.Async(), onKeyPressed)
}

func onKeyPressed(key int) {
    fmt.Println("Key pressed:", key)
}
```

### "Full" example

```go
package main

import (
	"fmt"
	"github.com/majfault/signal"
	"github.com/majfault/signal/dispatcher"
	"math/rand"
	"time"
)

type IntGenerator struct {
	IntGenerated signal.Signal1[int]
}

func (g *IntGenerator) start() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		g.IntGenerated.Emit(rand.Int())
	}
}

func onIntGenerated(i int) {
	fmt.Println("(func)   int generated:", i)
}

func main() {
	generator := &IntGenerator{
		IntGenerated: signal.Signal1[int]{},
	}

	// connects a declared function using the Async dispatcher
	generator.IntGenerated.Connect(dispatcher.Async(), onIntGenerated)

	// creates a Queued dispatcher with a capacity of 8
	queuedDispatcher := &dispatcher.Queued{Cap: 8}
	// connects a lambda to using declared Queued dispatcher
	slot := generator.IntGenerated.Connect(queuedDispatcher, func(i int) {
		fmt.Println("(lambda) int generated:", i)
	})

	go generator.start()

	<-time.After(3 * time.Second)
	// disconnects slot from signal after 3s
	generator.IntGenerated.Disconnect(slot)

	<-time.After(3 * time.Second)
}
```

Output:

    (lambda) int generated: 5577006791947779410
    (func)   int generated: 5577006791947779410
    (lambda) int generated: 8674665223082153551
    (func)   int generated: 8674665223082153551
    (lambda) int generated: 6129484611666145821
    (func)   int generated: 6129484611666145821
    (func)   int generated: 4037200794235010051
    (func)   int generated: 3916589616287113937
    (func)   int generated: 6334824724549167320