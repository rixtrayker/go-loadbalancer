package v1

import (
	"net/http"

	"github.com/rixtrayker/go-loadbalancer/internal/admin"
)

// RegisterAdminAPI registers all admin API endpoints
func RegisterAdminAPI(mux *http.ServeMux) {
	// Register admin API routes
	mux.HandleFunc("/api/v1/admin/status", admin.StatusHandler)
	mux.HandleFunc("/api/v1/admin/backends", admin.BackendsHandler)
	mux.HandleFunc("/api/v1/admin/metrics", admin.MetricsHandler)
	// Add more admin API endpoints as needed
}
