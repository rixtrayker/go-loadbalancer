package algorithms

import (
	"net/http"
	"sync/atomic"

	"github.com/rixtrayker/go-loadbalancer/internal/backend"
)

// RoundRobin implements the round-robin load balancing algorithm
type RoundRobin struct {
	backends []*backend.Backend
	current  uint32
}

// NewRoundRobin creates a new round-robin algorithm instance
func NewRoundRobin(backends []*backend.Backend) *RoundRobin {
	return &RoundRobin{
		backends: backends,
		current:  0,
	}
}

// NextBackend selects the next backend in a round-robin fashion
func (rr *RoundRobin) NextBackend(r *http.Request) *backend.Backend {
	if len(rr.backends) == 0 {
		return nil
	}

	// Get only healthy backends
	healthyBackends := make([]*backend.Backend, 0, len(rr.backends))
	for _, b := range rr.backends {
		if b.IsHealthy() {
			healthyBackends = append(healthyBackends, b)
		}
	}

	if len(healthyBackends) == 0 {
		return nil
	}

	// Get the next index in a thread-safe way
	idx := int(atomic.AddUint32(&rr.current, 1) - 1) % len(healthyBackends)
	return healthyBackends[idx]
}
