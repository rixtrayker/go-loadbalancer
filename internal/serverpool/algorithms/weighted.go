package algorithms

import (
	"net/http"
	"sync"

	"github.com/rixtrayker/go-loadbalancer/internal/backend"
)

// Weighted implements the weighted round-robin load balancing algorithm
type Weighted struct {
	backends []*backend.Backend
	current  int
	weights  []int
	gcd      int
	maxW     int
	cw       int
	mutex    sync.Mutex
}

// NewWeighted creates a new weighted round-robin algorithm instance
func NewWeighted(backends []*backend.Backend) *Weighted {
	weights := make([]int, len(backends))
	maxW := 0

	for i, b := range backends {
		weights[i] = b.Weight
		if b.Weight > maxW {
			maxW = b.Weight
		}
	}

	return &Weighted{
		backends: backends,
		weights:  weights,
		maxW:     maxW,
		gcd:      gcdOfArray(weights),
	}
}

// NextBackend selects the next backend using weighted round-robin
func (w *Weighted) NextBackend(r *http.Request) *backend.Backend {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if len(w.backends) == 0 {
		return nil
	}

	// Get only healthy backends
	healthyBackends := make([]*backend.Backend, 0, len(w.backends))
	healthyWeights := make([]int, 0, len(w.backends))

	for i, b := range w.backends {
		if b.IsHealthy() {
			healthyBackends = append(healthyBackends, b)
			healthyWeights = append(healthyWeights, w.weights[i])
		}
	}

	if len(healthyBackends) == 0 {
		return nil
	}

	// If only one backend, return it
	if len(healthyBackends) == 1 {
		return healthyBackends[0]
	}

	// Weighted round-robin algorithm
	for {
		w.current = (w.current + 1) % len(healthyBackends)
		if w.current == 0 {
			w.cw = w.cw - w.gcd
			if w.cw <= 0 {
				w.cw = w.maxW
				if w.cw == 0 {
					return healthyBackends[0]
				}
			}
		}

		if healthyWeights[w.current] >= w.cw {
			return healthyBackends[w.current]
		}
	}
}

// gcd calculates the greatest common divisor of two integers
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// gcdOfArray calculates the greatest common divisor of an array of integers
func gcdOfArray(arr []int) int {
	if len(arr) == 0 {
		return 1
	}
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] != 0 {
			result = gcd(result, arr[i])
		}
	}
	if result == 0 {
		return 1
	}
	return result
}
