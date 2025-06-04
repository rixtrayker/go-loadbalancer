package algorithms

import (
	"github.com/amr/go-loadbalancer/internal/backend"
)

// Algorithm defines the interface for load balancing algorithms
type Algorithm interface {
	// Next returns the next backend to use
	Next(backends []*backend.Backend) *backend.Backend
	// Name returns the name of the algorithm
	Name() string
	// AddBackend adds a backend to the algorithm
	AddBackend(*backend.Backend)
	// RemoveBackend removes a backend from the algorithm
	RemoveBackend(string)
	// GetBackends returns all backends in the algorithm
	GetBackends() []*backend.Backend
	// GetAvailableBackends returns all available backends
	GetAvailableBackends() []*backend.Backend
	// Size returns the number of backends
	Size() int
}
