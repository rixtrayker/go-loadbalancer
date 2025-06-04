package serverpool

import (
	"errors"
	"net/http"
	"sync"

	"github.com/rixtrayker/go-loadbalancer/configs"
	"github.com/rixtrayker/go-loadbalancer/internal/backend"
	"github.com/rixtrayker/go-loadbalancer/internal/serverpool/algorithms"
)

// Pool represents a group of backend servers
type Pool struct {
	Name      string
	Backends  []*backend.Backend
	Algorithm algorithms.Algorithm
	mutex     sync.RWMutex
}

// NewPool creates a new backend pool
func NewPool(config configs.BackendPoolConfig) (*Pool, error) {
	if len(config.Backends) == 0 {
		return nil, errors.New("no backends provided")
	}

	// Create backends
	backends := make([]*backend.Backend, 0, len(config.Backends))
	for _, backendConfig := range config.Backends {
		b, err := backend.NewBackend(backendConfig.URL, backendConfig.Weight)
		if err != nil {
			return nil, err
		}
		backends = append(backends, b)
	}

	// Create load balancing algorithm
	var algorithm algorithms.Algorithm
	switch config.Algorithm {
	case "round_robin":
		algorithm = algorithms.NewRoundRobin(backends)
	case "least_conn":
		algorithm = algorithms.NewLeastConn(backends)
	case "weighted":
		algorithm = algorithms.NewWeighted(backends)
	default:
		algorithm = algorithms.NewRoundRobin(backends)
	}

	return &Pool{
		Name:      config.Name,
		Backends:  backends,
		Algorithm: algorithm,
	}, nil
}

// NextBackend selects the next backend for a request
func (p *Pool) NextBackend(r *http.Request) (*backend.Backend, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Get healthy backends
	healthyBackends := make([]*backend.Backend, 0, len(p.Backends))
	for _, b := range p.Backends {
		if b.IsHealthy() {
			healthyBackends = append(healthyBackends, b)
		}
	}

	if len(healthyBackends) == 0 {
		return nil, errors.New("no healthy backends available")
	}

	// Select backend using the algorithm
	b := p.Algorithm.NextBackend(r)
	if b == nil {
		return nil, errors.New("failed to select backend")
	}

	// Update backend stats
	b.IncrementRequests()
	b.IncrementConnections()

	return b, nil
}

// MarkBackendStatus updates the health status of a backend
func (p *Pool) MarkBackendStatus(url string, healthy bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, b := range p.Backends {
		if b.URL.String() == url {
			b.SetHealth(healthy)
			break
		}
	}
}
