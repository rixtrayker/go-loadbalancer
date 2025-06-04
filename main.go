package main

import (
	"flag"
	"log"

	"github.com/rixtrayker/go-loadbalancer/internal/app"
	"github.com/rixtrayker/go-loadbalancer/pkg/logging"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "configs/config.yml", "Path to configuration file")
	flag.Parse()

	// Initialize logger
	logger := logging.NewLogger()
	logger.Info("Starting Go Load Balancer")

	// Create and run application
	application, err := app.New(*configPath)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
