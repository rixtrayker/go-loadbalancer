package http

import (
	"io"
	"net/http"
	"time"

	"github.com/amr/go-loadbalancer/internal/policy"
	"github.com/amr/go-loadbalancer/internal/routing"
	"go.uber.org/zap"
)

// Handler implements the HTTP handler for the load balancer
type Handler struct {
	router      *routing.Router
	policyChain *policy.PolicyChain
	logger      *zap.Logger
	client      *http.Client
}

// NewHandler creates a new HTTP handler
func NewHandler(router *routing.Router, policyChain *policy.PolicyChain, logger *zap.Logger) *Handler {
	return &Handler{
		router:      router,
		policyChain: policyChain,
		logger:      logger,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Apply request policies
	if err := h.policyChain.Apply(r, nil); err != nil {
		h.logger.Error("Policy application failed", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Route request to backend
	backend := h.router.RouteRequest(r)
	if backend == nil {
		h.logger.Error("No available backend")
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	// Create proxy request
	backendURL := backend.URL.String()
	proxyReq, err := http.NewRequest(r.Method, backendURL+r.URL.Path, r.Body)
	if err != nil {
		h.logger.Error("Failed to create proxy request", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Forward request to backend
	backend.IncrementConn()
	defer backend.DecrementConn()

	resp, err := h.client.Do(proxyReq)
	if err != nil {
		h.logger.Error("Failed to forward request", zap.Error(err))
		backend.RecordFailure()
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Apply response policies
	if err := h.policyChain.Apply(r, resp); err != nil {
		h.logger.Error("Response policy application failed", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set response status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	if _, err := io.Copy(w, resp.Body); err != nil {
		h.logger.Error("Failed to copy response body", zap.Error(err))
		return
	}

	// Record success
	backend.RecordSuccess()
}
