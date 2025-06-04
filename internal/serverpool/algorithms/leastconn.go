package algorithms

import (
	"net/http"

	"github.com/rixtrayker/go-loadbalancer/internal/backend"
)

// LeastConn implements the least connections load balancing algorithm
type LeastConn struct {
	backends []*backend.Backend
}

// NewLeastConn creates a new least connections algorithm instance
func NewLeastConn(backends []*backend.Backend) *LeastConn {
	return &LeastConn{
		backends: backends,
	}
}

// NextBackend selects the backend with the least active connections
func (lc *LeastConn) NextBackend(r *http.Request) *backend.Backend {
	if len(lc.backends) == 0 {
		return nil
	}

	// Get only healthy backends
	healthyBackends := make([]*backend.Backend, 0, len(lc.backends))
	for _, b := range lc.backends {
		if b.IsHealthy() {
			healthyBackends = append(healthyBackends, b)
		}
	}

	if len(healthyBackends) == 0 {
		return nil
	}

	// Find the backend with the least connections
	var selected *backend.Backend
	minConn := -1

	for _, b := range healthyBackends {
		conns := b.GetActiveConnections()
		if minConn == -1 || conns < minConn {
			minConn = conns
			selected = b
		}
	}

	return selected
}
