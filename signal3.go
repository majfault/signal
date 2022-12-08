package signal

import (
	"github.com/majfault/signal/dispatcher"
	"sort"
	"sync"
)

type Slot3[T, U, V any] struct {
	dispatcher dispatcher.Interface
	call       func(T, U, V)
}

type Signal3[T, U, V any] struct {
	mu    sync.RWMutex
	slots []*Slot3[T, U, V]
}

func (s *Signal3[T, U, V]) Connect(dispatcher dispatcher.Interface, slot func(t T, u U, v V)) *Slot3[T, U, V] {
	s.mu.Lock()
	defer s.mu.Unlock()
	ps := &Slot3[T, U, V]{dispatcher, slot}
	s.slots = append(s.slots, ps)
	sort.Slice(s.slots, func(i, j int) bool { return s.slots[i].dispatcher.Priority() > s.slots[j].dispatcher.Priority() })
	return ps
}

func (s *Signal3[T, U, V]) Disconnect(slot *Slot3[T, U, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots = deleteFirst(s.slots, slot)
}

func (s *Signal3[T, U, V]) Emit(t T, u U, v V) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.slots {
		slot := s.slots[i]
		slot.dispatcher.Dispatch(func() { slot.call(t, u, v) })
	}
}

func (s *Signal3[T, U, V]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots = nil
	return nil
}
