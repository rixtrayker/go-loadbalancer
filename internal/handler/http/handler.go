package http

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gorilla/mux"
	"github.com/rixtrayker/go-loadbalancer/configs"
	"github.com/rixtrayker/go-loadbalancer/internal/policy"
	"github.com/rixtrayker/go-loadbalancer/internal/serverpool"
	"github.com/rixtrayker/go-loadbalancer/pkg/logging"
	"github.com/rixtrayker/go-loadbalancer/pkg/metrics"
)

// Handler handles HTTP requests
type Handler struct {
	router  *mux.Router
	config  *configs.Config
	logger  *logging.Logger
	metrics *metrics.Metrics
}

// NewHandler creates a new HTTP handler
func NewHandler(config *configs.Config, logger *logging.Logger, metrics *metrics.Metrics) *Handler {
	h := &Handler{
		router:  mux.NewRouter(),
		config:  config,
		logger:  logger,
		metrics: metrics,
	}

	h.setupRoutes()
	return h
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	h.router.ServeHTTP(w, r)
	duration := time.Since(start)
	h.metrics.RecordRequestDuration(r.URL.Path, r.Method, duration)
}

// setupRoutes configures the HTTP routes
func (h *Handler) setupRoutes() {
	// Setup backend pools
	pools := make(map[string]*serverpool.Pool)
	for _, poolConfig := range h.config.BackendPools {
		pool, err := serverpool.NewPool(poolConfig)
		if err != nil {
			h.logger.Error("Failed to create backend pool", "name", poolConfig.Name, "error", err)
			continue
		}
		pools[poolConfig.Name] = pool
	}

	// Setup routing rules
	for _, rule := range h.config.RoutingRules {
		pool, ok := pools[rule.TargetPool]
		if !ok {
			h.logger.Error("Target pool not found", "pool", rule.TargetPool)
			continue
		}

		// Create route handler
		handler := func(w http.ResponseWriter, r *http.Request) {
			// Apply policies
			for _, policyConfig := range rule.Policies {
				if err := policy.Apply(policyConfig, r); err != nil {
					h.logger.Error("Policy application failed", "error", err)
					http.Error(w, "Policy violation", http.StatusForbidden)
					return
				}
			}

			// Select backend
			backend, err := pool.NextBackend(r)
			if err != nil {
				h.logger.Error("Failed to select backend", "error", err)
				http.Error(w, "No backend available", http.StatusServiceUnavailable)
				return
			}

			// Proxy the request
			proxy := httputil.NewSingleHostReverseProxy(backend.URL)
			proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
				h.logger.Error("Proxy error", "error", err)
				http.Error(w, "Backend error", http.StatusBadGateway)
			}

			h.logger.Info("Proxying request", "path", r.URL.Path, "backend", backend.URL.String())
			proxy.ServeHTTP(w, r)
		}

		// Register route
		route := h.router.NewRoute()
		if rule.Match.Host != "" {
			route = route.Host(rule.Match.Host)
		}
		if rule.Match.Path != "" {
			route = route.Path(rule.Match.Path)
		}
		if rule.Match.Method != "" {
			route = route.Methods(rule.Match.Method)
		}
		for k, v := range rule.Match.Headers {
			route = route.HeadersRegexp(k, v)
		}

		route.HandlerFunc(handler)
		h.logger.Info("Registered route", "host", rule.Match.Host, "path", rule.Match.Path)
	}

	// Add catch-all route
	h.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("No matching route", "path", r.URL.Path)
		http.Error(w, "No matching route", http.StatusNotFound)
	})
}
