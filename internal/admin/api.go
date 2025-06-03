package admin

import (
	"encoding/json"
	"net/http"

	"github.com/amr/go-loadbalancer/internal/backend"
	"github.com/amr/go-loadbalancer/internal/serverpool"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// API implements the admin API
type API struct {
	pool   *serverpool.Pool
	logger *logrus.Logger
}

// NewAPI creates a new admin API
func NewAPI(pool *serverpool.Pool, logger *logrus.Logger) *API {
	return &API{
		pool:   pool,
		logger: logger,
	}
}

// RegisterRoutes registers the admin API routes
func (a *API) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/admin/backends", a.listBackends).Methods("GET")
	r.HandleFunc("/admin/backends", a.addBackend).Methods("POST")
	r.HandleFunc("/admin/backends/{url}", a.removeBackend).Methods("DELETE")
	r.HandleFunc("/admin/health", a.healthCheck).Methods("GET")
}

// listBackends returns a list of all backends
func (a *API) listBackends(w http.ResponseWriter, r *http.Request) {
	backends := a.pool.GetBackends()
	json.NewEncoder(w).Encode(backends)
}

// addBackend adds a new backend
func (a *API) addBackend(w http.ResponseWriter, r *http.Request) {
	var b backend.Backend
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.pool.AddBackend(&b)
	w.WriteHeader(http.StatusCreated)
}

// removeBackend removes a backend
func (a *API) removeBackend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["url"]

	a.pool.RemoveBackend(url)
	w.WriteHeader(http.StatusNoContent)
}

// healthCheck returns the health status of the load balancer
func (a *API) healthCheck(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status    string `json:"status"`
		Backends  int    `json:"backends"`
		Available int    `json:"available"`
	}{
		Status:    "healthy",
		Backends:  a.pool.Size(),
		Available: len(a.pool.GetAvailableBackends()),
	}

	json.NewEncoder(w).Encode(status)
} 