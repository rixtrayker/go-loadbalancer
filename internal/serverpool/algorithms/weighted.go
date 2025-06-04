package algorithms

import (
	"sync/atomic"

	"github.com/amr/go-loadbalancer/internal/backend"
)

// Weighted implements the weighted round-robin load balancing algorithm
type Weighted struct {
	BaseAlgorithm
	current uint64
}

// NewWeighted creates a new weighted round-robin algorithm
func NewWeighted() *Weighted {
	return &Weighted{
		BaseAlgorithm: BaseAlgorithm{
			name: "weighted-round-robin",
		},
	}
}

// Next returns the next backend using weighted round-robin selection
func (w *Weighted) Next(backends []*backend.Backend) *backend.Backend {
	if len(backends) == 0 {
		return nil
	}

	// Calculate total weight
	var totalWeight int
	for _, b := range backends {
		if b.IsAvailable() {
			totalWeight += b.Weight
		}
	}

	if totalWeight == 0 {
		return nil
	}

	// Get the next weight index
	next := atomic.AddUint64(&w.current, 1)
	weightIndex := int(next % uint64(totalWeight))

	// Find the backend corresponding to the weight index
	var currentWeight int
	for _, b := range backends {
		if !b.IsAvailable() {
			continue
		}

		currentWeight += b.Weight
		if weightIndex < currentWeight {
			return b
		}
	}

	return nil
}

// AddBackend adds a backend to the algorithm
func (w *Weighted) AddBackend(b *backend.Backend) {
	// Implementation needed
}

// RemoveBackend removes a backend from the algorithm
func (w *Weighted) RemoveBackend(name string) {
	// Implementation needed
}

// GetBackends returns all backends
func (w *Weighted) GetBackends() []*backend.Backend {
	// Implementation needed
	return nil
}

// GetAvailableBackends returns all available backends
func (w *Weighted) GetAvailableBackends() []*backend.Backend {
	// Implementation needed
	return nil
}

// Size returns the number of backends
func (w *Weighted) Size() int {
	// Implementation needed
	return 0
}