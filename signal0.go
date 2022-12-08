package signal

import (
	"github.com/majfault/signal/dispatcher"
	"sort"
	"sync"
)

type Slot0 struct {
	dispatcher dispatcher.Interface
	call       func()
}

type Signal0 struct {
	mu    sync.RWMutex
	slots []*Slot0
}

func (s *Signal0) Connect(dispatcher dispatcher.Interface, slot func()) *Slot0 {
	s.mu.Lock()
	defer s.mu.Unlock()
	ps := &Slot0{dispatcher, slot}
	s.slots = append(s.slots, ps)
	sort.Slice(s.slots, func(i, j int) bool { return s.slots[i].dispatcher.Priority() > s.slots[j].dispatcher.Priority() })
	return ps
}

func (s *Signal0) Disconnect(slot *Slot0) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots = deleteFirst(s.slots, slot)
}

func (s *Signal0) Emit() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.slots {
		slot := s.slots[i]
		slot.dispatcher.Dispatch(slot.call)
	}
}

func (s *Signal0) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots = nil
	return nil
}
