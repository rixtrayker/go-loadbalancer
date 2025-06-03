package serverpool

import (
	"sync"

	"github.com/amr/go-loadbalancer/internal/backend"
)

// Pool represents a pool of backend servers
type Pool struct {
	backends []*backend.Backend
	current  uint64
	mu       sync.RWMutex
}

// NewPool creates a new server pool
func NewPool() *Pool {
	return &Pool{
		backends: make([]*backend.Backend, 0),
	}
}

// AddBackend adds a backend to the pool
func (p *Pool) AddBackend(b *backend.Backend) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.backends = append(p.backends, b)
}

// RemoveBackend removes a backend from the pool
func (p *Pool) RemoveBackend(url string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, b := range p.backends {
		if b.URL == url {
			p.backends = append(p.backends[:i], p.backends[i+1:]...)
			return
		}
	}
}

// GetNextBackend returns the next available backend
func (p *Pool) GetNextBackend() *backend.Backend {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.backends) == 0 {
		return nil
	}

	// Try to find an available backend
	for i := 0; i < len(p.backends); i++ {
		backend := p.backends[i]
		if backend.IsAvailable() {
			return backend
		}
	}

	return nil
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