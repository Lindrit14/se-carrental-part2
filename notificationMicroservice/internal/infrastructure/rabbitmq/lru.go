package rabbitmq

import "sync"

// lruSeen is a thread-safe in-memory set with a fixed capacity.
// When full, the oldest entry is evicted (FIFO order via a ring buffer).
type lruSeen struct {
	mu   sync.Mutex
	set  map[string]struct{}
	ring []string
	cap  int
	pos  int
}

func newLRUSeen(capacity int) *lruSeen {
	return &lruSeen{
		set:  make(map[string]struct{}, capacity),
		ring: make([]string, capacity),
		cap:  capacity,
	}
}

func (s *lruSeen) has(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.set[id]
	return ok
}

func (s *lruSeen) add(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.set[id]; ok {
		return
	}
	// Evict oldest if at capacity.
	if evict := s.ring[s.pos]; evict != "" {
		delete(s.set, evict)
	}
	s.set[id] = struct{}{}
	s.ring[s.pos] = id
	s.pos = (s.pos + 1) % s.cap
}
