package serverpool

import (
	"sync"

	"github.com/amr/go-loadbalancer/internal/backend"
	"github.com/amr/go-loadbalancer/internal/serverpool/algorithms"
)

// Pool represents a pool of backend servers
type Pool struct {
	backends  []*backend.Backend
	algorithm algorithms.Algorithm
	current  uint64
	mu       sync.RWMutex
}

// NewPool creates a new server pool
func NewPool() *Pool {
	return &Pool{
		backends: make([]*backend.Backend, 0),
	}
}

// NewBackend creates a new backend instance
func NewBackend(url string, weight int) (*backend.Backend, error) {
	return backend.NewBackend("", url, weight, 100)
}

// AddBackend adds a backend to the pool
func (p *Pool) AddBackend(b *backend.Backend) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.backends = append(p.backends, b)
}

// RemoveBackend removes a backend from the pool
func (p *Pool) RemoveBackend(urlStr string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, b := range p.backends {
		if b.URL.String() == urlStr {
			p.backends = append(p.backends[:i], p.backends[i+1:]...)
			return
		}
	}
}

// GetBackends returns all backends in the pool
func (p *Pool) GetBackends() []*backend.Backend {
	p.mu.RLock()
	defer p.mu.RUnlock()

	backends := make([]*backend.Backend, len(p.backends))
	copy(backends, p.backends)
	return backends
}

// GetAvailableBackends returns all available backends
func (p *Pool) GetAvailableBackends() []*backend.Backend {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var available []*backend.Backend
	for _, b := range p.backends {
		if b.IsAvailable() {
			available = append(available, b)
		}
	}
	return available
}

// Size returns the number of backends in the pool
func (p *Pool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.backends)
}

// SetAlgorithm sets the load balancing algorithm
func (p *Pool) SetAlgorithm(algorithm algorithms.Algorithm) {
	p.algorithm = algorithm
}

// GetNextBackend returns the next backend using the configured algorithm
func (p *Pool) GetNextBackend() *backend.Backend {
	if p.algorithm == nil {
		return nil
	}
	return p.algorithm.Next(p.GetAvailableBackends())
}
