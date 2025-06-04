package algorithms

import (
	"sync"

	"github.com/amr/go-loadbalancer/internal/backend"
)

// BaseAlgorithm provides common functionality for algorithms
type BaseAlgorithm struct {
	name     string
	backends []*backend.Backend
	mu       sync.RWMutex
}

// Name returns the name of the algorithm
func (b *BaseAlgorithm) Name() string {
	return b.name
}

// AddBackend adds a backend to the algorithm
func (b *BaseAlgorithm) AddBackend(backend *backend.Backend) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.backends = append(b.backends, backend)
}

// RemoveBackend removes a backend from the algorithm
func (b *BaseAlgorithm) RemoveBackend(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, backend := range b.backends {
		if backend.Name == name {
			b.backends = append(b.backends[:i], b.backends[i+1:]...)
			return
		}
	}
}

// GetBackends returns all backends
func (b *BaseAlgorithm) GetBackends() []*backend.Backend {
	b.mu.RLock()
	defer b.mu.RUnlock()
	result := make([]*backend.Backend, len(b.backends))
	copy(result, b.backends)
	return result
}

// GetAvailableBackends returns all available backends
func (b *BaseAlgorithm) GetAvailableBackends() []*backend.Backend {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	available := make([]*backend.Backend, 0)
	for _, backend := range b.backends {
		if backend.IsAvailable() {
			available = append(available, backend)
		}
	}
	return available
}

// Size returns the number of backends
func (b *BaseAlgorithm) Size() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.backends)
}
