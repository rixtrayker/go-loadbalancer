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
