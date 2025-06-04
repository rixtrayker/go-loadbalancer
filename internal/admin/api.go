package admin

import (
	"encoding/json"
	"net/http"

	"github.com/rixtrayker/go-loadbalancer/internal/logging"
	"github.com/rixtrayker/go-loadbalancer/internal/serverpool"
)

// API handles admin API requests
type API struct {
	pools   map[string]*serverpool.Pool
	logger  *logging.Logger
}

// NewAPI creates a new admin API
func NewAPI(
	pools map[string]*serverpool.Pool,
	logger *logging.Logger,
) *API {
	return &API{
		pools:   pools,
		logger:  logger,
	}
}

// RegisterHandlers registers admin API handlers
func (a *API) RegisterHandlers(mux *http.ServeMux, basePath string) {
	mux.HandleFunc(basePath+"/status", a.handleStatus)
	mux.HandleFunc(basePath+"/backends", a.handleBackends)
	mux.HandleFunc(basePath+"/metrics", a.handleMetrics)
}

// handleStatus handles status requests
func (a *API) handleStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status": "ok",
		"pools":  len(a.pools),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// handleBackends handles backend management requests
func (a *API) handleBackends(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// List backends
		result := make(map[string][]map[string]interface{})

		for name, pool := range a.pools {
			backends := make([]map[string]interface{}, 0, len(pool.Backends))
			for _, b := range pool.Backends {
				backends = append(backends, map[string]interface{}{
					"url":            b.URL.String(),
					"healthy":        b.IsHealthy(),
					"active_conns":   b.GetActiveConnections(),
					"total_requests": b.GetTotalRequests(),
					"weight":         b.Weight,
				})
			}
			result[name] = backends
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

	case http.MethodPost:
		// Update backend status
		var req struct {
			Pool    string `json:"pool"`
			URL     string `json:"url"`
			Healthy bool   `json:"healthy"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		pool, ok := a.pools[req.Pool]
		if !ok {
			http.Error(w, "Pool not found", http.StatusNotFound)
			return
		}

		pool.MarkBackendStatus(req.URL, req.Healthy)
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleMetrics handles metrics requests
func (a *API) handleMetrics(w http.ResponseWriter, r *http.Request) {
	// This would typically use Prometheus HTTP handler
	// For now, just return a simple status
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "metrics available at /metrics endpoint",
	})
}

// StatusHandler handles status requests
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// BackendsHandler handles backend management requests
func BackendsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// MetricsHandler handles metrics requests
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
