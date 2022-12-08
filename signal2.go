package signal

import (
	"github.com/majfault/signal/dispatcher"
	"sort"
	"sync"
)

type Slot2[T, U any] struct {
	dispatcher dispatcher.Interface
	call       func(T, U)
}

type Signal2[T, U any] struct {
	mu    sync.RWMutex
	slots []*Slot2[T, U]
}

func (s *Signal2[T, U]) Connect(dispatcher dispatcher.Interface, slot func(t T, u U)) *Slot2[T, U] {
	s.mu.Lock()
	defer s.mu.Unlock()
	ps := &Slot2[T, U]{dispatcher, slot}
	s.slots = append(s.slots, ps)
	sort.Slice(s.slots, func(i, j int) bool { return s.slots[i].dispatcher.Priority() > s.slots[j].dispatcher.Priority() })
	return ps
}

func (s *Signal2[T, U]) Disconnect(slot *Slot2[T, U]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots = deleteFirst(s.slots, slot)
}

func (s *Signal2[T, U]) Emit(t T, u U) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.slots {
		slot := s.slots[i]
		slot.dispatcher.Dispatch(func() { slot.call(t, u) })
	}
}

func (s *Signal2[T, U]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots = nil
	return nil
}
