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
