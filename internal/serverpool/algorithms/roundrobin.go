package algorithms

import (
	"sync/atomic"

	"github.com/amr/go-loadbalancer/internal/backend"
)

// RoundRobin implements the round-robin load balancing algorithm
type RoundRobin struct {
	BaseAlgorithm
	current uint64
}

// NewRoundRobin creates a new round-robin algorithm
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		BaseAlgorithm: BaseAlgorithm{
			name: "round-robin",
		},
	}
}

// Next returns the next backend using round-robin selection
func (rr *RoundRobin) Next(backends []*backend.Backend) *backend.Backend {
	if len(backends) == 0 {
		return nil
	}

	// Try to find an available backend
	for i := 0; i < len(backends); i++ {
		// Get the next backend index using atomic operations
		next := atomic.AddUint64(&rr.current, 1)
		index := next % uint64(len(backends))
		backend := backends[index]

		if backend.IsAvailable() {
			return backend
		}
	}

	return nil
}

// AddBackend adds a backend to the algorithm
func (rr *RoundRobin) AddBackend(b *backend.Backend) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	rr.backends = append(rr.backends, b)
}

// RemoveBackend removes a backend from the algorithm
func (rr *RoundRobin) RemoveBackend(name string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for i, b := range rr.backends {
		if b.Name == name {
			rr.backends = append(rr.backends[:i], rr.backends[i+1:]...)
			if rr.next >= len(rr.backends) {
				rr.next = 0
			}
			return
		}
	}
}

// GetBackends returns all backends
func (rr *RoundRobin) GetBackends() []*backend.Backend {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	return rr.backends
}

// GetAvailableBackends returns all available backends
func (rr *RoundRobin) GetAvailableBackends() []*backend.Backend {
	rr.mu.RLock()
	defer rr.mu.RUnlock()

	available := make([]*backend.Backend, 0)
	for _, b := range rr.backends {
		if b.IsAvailable() {
			available = append(available, b)
		}
	}
	return available
}

// Size returns the number of backends
func (rr *RoundRobin) Size() int {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	return len(rr.backends)
} 