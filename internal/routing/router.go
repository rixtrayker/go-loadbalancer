package routing

import (
	"net/http"
	"sync"

	"github.com/amr/go-loadbalancer/internal/backend"
	"github.com/amr/go-loadbalancer/internal/serverpool"
)

// Router handles request routing to backends
type Router struct {
	pool     *serverpool.Pool
	rules    []*Rule
	mu       sync.RWMutex
}

// NewRouter creates a new router
func NewRouter(pool *serverpool.Pool) *Router {
	return &Router{
		pool:  pool,
		rules: make([]*Rule, 0),
	}
}

// AddRule adds a routing rule
func (r *Router) AddRule(rule *Rule) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules = append(r.rules, rule)
}

// RemoveRule removes a routing rule
func (r *Router) RemoveRule(rule *Rule) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, existing := range r.rules {
		if existing == rule {
			r.rules = append(r.rules[:i], r.rules[i+1:]...)
			return
		}
	}
}

// RouteRequest routes an incoming request to the appropriate backend
func (r *Router) RouteRequest(req *http.Request) *backend.Backend {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// First try to match against rules
	for _, rule := range r.rules {
		if rule.Matches(req) {
			return r.pool.GetNextBackend()
		}
	}

	// If no rules match, use the default backend
	return r.pool.GetNextBackend()
}

// GetRules returns all routing rules
func (r *Router) GetRules() []*Rule {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rules := make([]*Rule, len(r.rules))
	copy(rules, r.rules)
	return rules
} 