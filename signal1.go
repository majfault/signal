package signal

import (
	"github.com/majfault/signal/dispatcher"
	"sort"
	"sync"
)

type Slot1[T any] struct {
	dispatcher dispatcher.Interface
	call       func(T)
}

type Signal1[T any] struct {
	mu    sync.RWMutex
	slots []*Slot1[T]
}

func (s *Signal1[T]) Connect(dispatcher dispatcher.Interface, slot func(t T)) *Slot1[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	ps := &Slot1[T]{dispatcher, slot}
	s.slots = append(s.slots, ps)
	sort.Slice(s.slots, func(i, j int) bool { return s.slots[i].dispatcher.Priority() > s.slots[j].dispatcher.Priority() })
	return ps
}

func (s *Signal1[T]) Disconnect(slot *Slot1[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots = deleteFirst(s.slots, slot)
}

func (s *Signal1[T]) Emit(t T) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.slots {
		slot := s.slots[i]
		slot.dispatcher.Dispatch(func() { slot.call(t) })
	}
}

func (s *Signal1[T]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots = nil
	return nil
}
