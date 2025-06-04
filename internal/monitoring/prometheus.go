package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rixtrayker/go-loadbalancer/configs"
	"github.com/rixtrayker/go-loadbalancer/internal/logging"
)

// PrometheusServer represents the Prometheus metrics server
type PrometheusServer struct {
	server *http.Server
	logger *logging.Logger
}

// NewPrometheusServer creates a new Prometheus metrics server
func NewPrometheusServer(config configs.PrometheusConfig, logger *logging.Logger) *PrometheusServer {
	// Create a new mux for the metrics endpoint
	mux := http.NewServeMux()
	mux.Handle(config.Path, promhttp.Handler())

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
	}

	return &PrometheusServer{
		server: server,
		logger: logger,
	}
}

// Start starts the Prometheus metrics server
func (ps *PrometheusServer) Start() {
	go func() {
		ps.logger.Info("Starting Prometheus metrics server", "addr", ps.server.Addr)
		if err := ps.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ps.logger.Error("Prometheus server error", "error", err)
		}
	}()
}

// Stop stops the Prometheus metrics server
func (ps *PrometheusServer) Stop(ctx context.Context) error {
	ps.logger.Info("Stopping Prometheus metrics server")
	return ps.server.Shutdown(ctx)
}

// InitializePrometheus sets up the Prometheus metrics server
func InitializePrometheus(config configs.PrometheusConfig, logger *logging.Logger) (*PrometheusServer, error) {
	if !config.Enabled {
		logger.Info("Prometheus metrics server is disabled")
		return nil, nil
	}

	// Create and start the metrics server
	server := NewPrometheusServer(config, logger)
	server.Start()

	// Start metrics collector
	collector := NewMetricsCollector(logger)
	collector.Start(context.Background(), 15*time.Second)

	return server, nil
}
