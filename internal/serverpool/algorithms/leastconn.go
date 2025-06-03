package algorithms

import (
	"sync/atomic"

	"github.com/amr/go-loadbalancer/internal/backend"
)

// LeastConn implements the least connections load balancing algorithm
type LeastConn struct {
	BaseAlgorithm
}

// NewLeastConn creates a new least connections algorithm
func NewLeastConn() *LeastConn {
	return &LeastConn{
		BaseAlgorithm: BaseAlgorithm{
			name: "least-connections",
		},
	}
}

// Next returns the backend with the least number of active connections
func (lc *LeastConn) Next(backends []*backend.Backend) *backend.Backend {
	if len(backends) == 0 {
		return nil
	}

	var selected *backend.Backend
	var minConns int32 = -1

	for _, b := range backends {
		if !b.IsAvailable() {
			continue
		}

		conns := atomic.LoadInt32(&b.CurrentConns)
		if minConns == -1 || conns < minConns {
			minConns = conns
			selected = b
		}
	}

	return selected
}

// AddBackend adds a backend to the algorithm
func (lc *LeastConn) AddBackend(b *backend.Backend) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.backends = append(lc.backends, b)
}

// RemoveBackend removes a backend from the algorithm
func (lc *LeastConn) RemoveBackend(name string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	for i, b := range lc.backends {
		if b.Name == name {
			lc.backends = append(lc.backends[:i], lc.backends[i+1:]...)
			return
		}
	}
}

// GetBackends returns all backends
func (lc *LeastConn) GetBackends() []*backend.Backend {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.backends
}

// GetAvailableBackends returns all available backends
func (lc *LeastConn) GetAvailableBackends() []*backend.Backend {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	available := make([]*backend.Backend, 0)
	for _, b := range lc.backends {
		if b.IsAvailable() {
			available = append(available, b)
		}
	}
	return available
}

// Size returns the number of backends
func (lc *LeastConn) Size() int {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return len(lc.backends)
} 