package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rixtrayker/go-loadbalancer/configs"
	httpHandler "github.com/rixtrayker/go-loadbalancer/internal/handler/http"
	"github.com/rixtrayker/go-loadbalancer/internal/logging"
	"github.com/rixtrayker/go-loadbalancer/internal/middleware"
	"github.com/rixtrayker/go-loadbalancer/internal/monitoring"
	"github.com/rixtrayker/go-loadbalancer/internal/tracing"
)

// App represents the load balancer application
type App struct {
	config     *configs.Config
	httpServer *http.Server
	logger     *logging.Logger
	metrics    *monitoring.MetricsCollector
	tracer     *tracing.Tracer
}

// New creates a new application instance
func New(configPath string) (*App, error) {
	// Load configuration
	config, err := configs.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// Initialize logger with monitoring config
	logger := logging.NewLogger()
	if err := logger.Configure(config.Monitoring.Logging); err != nil {
		return nil, err
	}

	// Initialize metrics collector
	metricsCollector := monitoring.NewMetricsCollector(logger)
	if config.Monitoring.Prometheus.Enabled {
		if _, err := monitoring.InitializePrometheus(config.Monitoring.Prometheus, logger); err != nil {
			return nil, err
		}
	}

	// Initialize tracer if enabled
	var tracer *tracing.Tracer
	if config.Monitoring.Tracing.Enabled {
		tracer, err = tracing.NewTracer(config.Monitoring.Tracing, logger)
		if err != nil {
			return nil, err
		}
	}

	// Create the application
	app := &App{
		config:  config,
		logger:  logger,
		metrics: metricsCollector,
		tracer:  tracer,
	}

	// Setup HTTP server with monitoring middleware
	var handler http.Handler = httpHandler.NewHandler(config, logger)
	if config.Monitoring.Prometheus.Enabled {
		handler = middleware.MonitoringMiddleware(handler)
	}
	if config.Monitoring.Tracing.Enabled {
		handler = middleware.TracingMiddleware(handler)
	}

	app.httpServer = &http.Server{
		Addr:    config.Server.Address,
		Handler: handler,
	}

	return app, nil
}

// Run starts the application
func (a *App) Run() error {
	// Setup signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server
	go func() {
		a.logger.Info("Starting HTTP server on " + a.config.Server.Address)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("HTTP server error", "error", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	a.logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Error("Server forced to shutdown", "error", err)
		return err
	}

	a.logger.Info("Server gracefully stopped")
	return nil
}
