package routing

import (
	"net/http"

	"github.com/rixtrayker/go-loadbalancer/configs"
	"github.com/rixtrayker/go-loadbalancer/internal/serverpool"
	"github.com/rixtrayker/go-loadbalancer/pkg/logging"
)

// Router handles request routing
type Router struct {
	rules  []*Rule
	pools  map[string]*serverpool.Pool
	logger *logging.Logger
}

// NewRouter creates a new router
func NewRouter(
	config []configs.RoutingRuleConfig,
	pools map[string]*serverpool.Pool,
	logger *logging.Logger,
) *Router {
	router := &Router{
		rules:  make([]*Rule, 0, len(config)),
		pools:  pools,
		logger: logger,
	}

	// Create rules from config
	for _, ruleConfig := range config {
		rule := NewRule(ruleConfig)
		router.rules = append(router.rules, rule)
	}

	return router
}

// Route routes a request to the appropriate backend pool
func (r *Router) Route(req *http.Request) (*serverpool.Pool, error) {
	// Find matching rule
	for _, rule := range r.rules {
		if rule.Matches(req) {
			// Get the target pool
			pool, ok := r.pools[rule.TargetPool]
			if !ok {
				r.logger.Error("Target pool not found", "pool", rule.TargetPool)
				continue
			}
			return pool, nil
		}
	}

	return nil, ErrNoMatchingRule
}

// AddRule adds a new routing rule
func (r *Router) AddRule(rule *Rule) {
	r.rules = append(r.rules, rule)
}

// RemoveRule removes a routing rule
func (r *Router) RemoveRule(index int) {
	if index < 0 || index >= len(r.rules) {
		return
	}
	r.rules = append(r.rules[:index], r.rules[index+1:]...)
}
