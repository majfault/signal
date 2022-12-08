package dispatcher

// Interface describes the requirements for a slot dispatcher.
type Interface interface {
	Dispatch(slot func()) // dispatches slot by calling the function
	Priority() int        // priority used by Signal.Connect for slot call ordering, higher will be called before
}

// Asynchronous/goroutine-based dispatchers have higher priorities so that they are called before the direct dispatcher, which blocks.
const (
	AsyncPriority  = 50
	QueuedPriority = 20
	DirectPriority = 10
)
