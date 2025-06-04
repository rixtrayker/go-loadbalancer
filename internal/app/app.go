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
	"github.com/rixtrayker/go-loadbalancer/pkg/logging"
	"github.com/rixtrayker/go-loadbalancer/pkg/metrics"
)

// App represents the load balancer application
type App struct {
	config     *configs.Config
	httpServer *http.Server
	logger     *logging.Logger
	metrics    *metrics.Metrics
}

// New creates a new application instance
func New(configPath string) (*App, error) {
	// Load configuration
	config, err := configs.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// Initialize components
	logger := logging.NewLogger()
	metricsCollector := metrics.NewMetrics()

	// Create the application
	app := &App{
		config:  config,
		logger:  logger,
		metrics: metricsCollector,
	}

	// Setup HTTP server
	handler := httpHandler.NewHandler(config, logger, metricsCollector)
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
