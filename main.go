package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/amr/go-loadbalancer/config"
	"github.com/amr/go-loadbalancer/internal/app"
	"github.com/amr/go-loadbalancer/pkg/logging"
	"go.uber.org/zap"
)

const defaultConfigFile = "config/config.yml"

func main() {
	// Parse command line flags
	configFile := flag.String("config", defaultConfigFile, "Path to config file")
	flag.Parse()

	// Initialize logger
	logger := logging.DefaultLogger()
	defer func() {
		_ = logger.Sync()
	}()

	// Load configuration
	configLoader := config.NewLoader(*configFile)
	cfg, err := configLoader.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Create and start application
	application, err := app.New(cfg)
	if err != nil {
		logger.Fatal("Failed to create application", zap.Error(err))
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := application.Start(); err != nil {
			logger.Fatal("Failed to start application", zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	logger.Info("Received shutdown signal")

	if err := application.Shutdown(); err != nil {
		logger.Error("Error during shutdown", zap.Error(err))
	}
}
