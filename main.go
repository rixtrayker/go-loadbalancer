package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/amr/go-loadbalancer/internal/admin"
	"github.com/amr/go-loadbalancer/internal/handler/http"
	"github.com/amr/go-loadbalancer/internal/healthcheck"
	"github.com/amr/go-loadbalancer/internal/policy"
	"github.com/amr/go-loadbalancer/internal/policy/ratelimit"
	"github.com/amr/go-loadbalancer/internal/policy/security"
	"github.com/amr/go-loadbalancer/internal/policy/transform"
	"github.com/amr/go-loadbalancer/internal/routing"
	"github.com/amr/go-loadbalancer/internal/serverpool"
	"github.com/amr/go-loadbalancer/internal/serverpool/algorithms"
	"github.com/amr/go-loadbalancer/pkg/logging"
	"github.com/amr/go-loadbalancer/pkg/metrics"
	"github.com/amr/go-loadbalancer/pkg/tracer"
	"github.com/gorilla/mux"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8080", "Port to listen on")
	configFile := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	// Initialize logger
	logger := logging.DefaultLogger()

	// Initialize metrics
	metrics := metrics.NewMetrics()

	// Initialize tracer
	tracer := tracer.NewNoopTracer()

	// Initialize server pool
	pool := serverpool.NewPool()

	// Initialize load balancing algorithm
	algorithm := algorithms.NewRoundRobin()

	// Initialize health checker
	healthChecker := healthcheck.NewHealthChecker(30*time.Second, 5*time.Second, "http")
	healthChecker.Start()
	defer healthChecker.Stop()

	// Initialize router
	router := routing.NewRouter(pool)

	// Initialize policy chain
	policyChain := policy.NewPolicyChain()

	// Add rate limiting policy
	rateLimiter := ratelimit.NewRateLimiter(100, time.Minute)
	policyChain.AddPolicy(rateLimiter)

	// Add security policy
	acl := security.NewACL()
	policyChain.AddPolicy(acl)

	// Add transformation policy
	transformer := transform.NewTransformer()
	policyChain.AddPolicy(transformer)

	// Initialize HTTP handler
	handler := http.NewHandler(router, policyChain, logger)

	// Initialize admin API
	adminAPI := admin.NewAPI(pool, logger)

	// Create router
	r := mux.NewRouter()

	// Register admin routes
	adminAPI.RegisterRoutes(r)

	// Register main handler
	r.PathPrefix("/").Handler(handler)

	// Start server
	logger.Infof("Starting server on port %s", *port)
	if err := http.ListenAndServe(":"+*port, r); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
} 