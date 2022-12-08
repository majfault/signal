package signal

import (
	"github.com/majfault/signal/dispatcher"
	"sort"
	"sync"
)

type Slot4[T, U, V, W any] struct {
	dispatcher dispatcher.Interface
	call       func(T, U, V, W)
}

type Signal4[T, U, V, W any] struct {
	mu    sync.RWMutex
	slots []*Slot4[T, U, V, W]
}

func (s *Signal4[T, U, V, W]) Connect(dispatcher dispatcher.Interface, slot func(t T, u U, v V, w W)) *Slot4[T, U, V, W] {
	s.mu.Lock()
	defer s.mu.Unlock()
	ps := &Slot4[T, U, V, W]{dispatcher, slot}
	s.slots = append(s.slots, ps)
	sort.Slice(s.slots, func(i, j int) bool { return s.slots[i].dispatcher.Priority() > s.slots[j].dispatcher.Priority() })
	return ps
}

func (s *Signal4[T, U, V, W]) Disconnect(slot *Slot4[T, U, V, W]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots = deleteFirst(s.slots, slot)
}

func (s *Signal4[T, U, V, W]) Emit(t T, u U, v V, w W) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.slots {
		slot := s.slots[i]
		slot.dispatcher.Dispatch(func() { slot.call(t, u, v, w) })
	}
}

func (s *Signal4[T, U, V, W]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots = nil
	return nil
}
