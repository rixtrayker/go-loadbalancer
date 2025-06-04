package algorithms

import (
	"net/http"

	"github.com/rixtrayker/go-loadbalancer/internal/backend"
)

// Algorithm defines the interface for load balancing algorithms
type Algorithm interface {
	// NextBackend selects the next backend for a request
	NextBackend(r *http.Request) *backend.Backend
}
