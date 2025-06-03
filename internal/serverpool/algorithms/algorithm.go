package algorithms

import (
	"github.com/amr/go-loadbalancer/internal/backend"
)

// Algorithm defines the interface for load balancing algorithms
type Algorithm interface {
	// GetNextBackend returns the next backend to use
	GetNextBackend(backends []*backend.Backend) *backend.Backend
}

// RoundRobin implements round-robin load balancing
type RoundRobin struct {
	current int
}

// NewRoundRobin creates a new round-robin algorithm
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		current: 0,
	}
}

// GetNextBackend implements the Algorithm interface
func (rr *RoundRobin) GetNextBackend(backends []*backend.Backend) *backend.Backend {
	if len(backends) == 0 {
		return nil
	}

	backend := backends[rr.current]
	rr.current = (rr.current + 1) % len(backends)
	return backend
} 