package app

import (
	"net/http"
	"strconv"
	"time"

	"github.com/amr/go-loadbalancer/config"
	"github.com/amr/go-loadbalancer/internal/admin"
	lbhttp "github.com/amr/go-loadbalancer/internal/handler/http"
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
	"go.uber.org/zap"
)

// App represents the main application
type App struct {
	config      *config.Config
	logger      *zap.Logger
	metrics     *metrics.Metrics
	tracer      tracer.Tracer
	server      *http.Server
	pool        *serverpool.Pool
	router      *routing.Router
	policyChain *policy.PolicyChain
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	// Initialize logger
	logger := logging.DefaultLogger()

	// Initialize metrics
	metricsCollector := metrics.NewMetrics()

	// Initialize tracer
	var tracerInstance tracer.Tracer
	jaegerTracer, err := tracer.NewJaegerTracer("go-loadbalancer")
	if err != nil {
		logger.Warn("Failed to initialize Jaeger tracer, falling back to no-op tracer", zap.Error(err))
		tracerInstance = tracer.NewNoopTracer()
	} else {
		tracerInstance = jaegerTracer
	}

	// Initialize server pool
	pool := serverpool.NewPool()

	// Add backends from config
	for _, backend := range cfg.Backends {
		backend, err := serverpool.NewBackend(backend.URL, backend.Weight)
		if err != nil {
			logger.Error("Failed to create backend", 
				zap.String("url", backend.URL.String()),
				zap.Error(err))
			continue
		}
		pool.AddBackend(backend)
	}

	// Initialize load balancing algorithm based on config
	var algorithm algorithms.Algorithm
	switch cfg.Server.Algorithm {
	case "round-robin":
		algorithm = algorithms.NewRoundRobin()
	default:
		algorithm = algorithms.NewRoundRobin()
	}
	pool.SetAlgorithm(algorithm)

	// Initialize health checker
	healthChecker := healthcheck.NewHealthChecker(30*time.Second, 5*time.Second, "http")
	healthChecker.Start()

	// Initialize router
	router := routing.NewRouter(pool)

	// Initialize policy chain
	policyChain := policy.NewPolicyChain()

	// Add rate limiting policy if enabled
	if cfg.Policies.RateLimit.Enabled {
		rateLimiter := ratelimit.NewRateLimiter(
			cfg.Policies.RateLimit.RequestsPer,
			time.Duration(cfg.Policies.RateLimit.Period)*time.Second,
		)
		policyChain.AddPolicy(rateLimiter)
	}

	// Add security policy
	acl := security.NewACL()
	if len(cfg.Policies.Security.AllowedIPs) > 0 {
		if err := acl.SetAllowedIPs(cfg.Policies.Security.AllowedIPs); err != nil {
			logger.Warn("Failed to set allowed IPs", zap.Error(err))
		}
	}
	policyChain.AddPolicy(acl)

	// Add transformation policy
	transformer := transform.NewTransformer()
	if len(cfg.Policies.Transform.AddHeaders) > 0 {
		for k, v := range cfg.Policies.Transform.AddHeaders {
			transformer.AddHeaderTransformation(k, v)
		}
	}
	if len(cfg.Policies.Transform.RemoveHeaders) > 0 {
		for _, h := range cfg.Policies.Transform.RemoveHeaders {
			transformer.RemoveHeaderTransformation(h)
		}
	}
	policyChain.AddPolicy(transformer)

	// Initialize HTTP handler with zap logger
	handler := lbhttp.NewHandler(router, policyChain, logger)

	// Initialize admin API with zap logger
	adminAPI := admin.NewAPI(pool, logger)

	// Create router
	r := mux.NewRouter()

	// Register admin routes
	adminAPI.RegisterRoutes(r)

	// Register main handler
	r.PathPrefix("/").Handler(handler)

	// Create server with configured timeouts
	server := &http.Server{
		Addr:         cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	return &App{
		config:      cfg,
		logger:      logger,
		metrics:     metricsCollector,
		tracer:      tracerInstance,
		server:      server,
		pool:        pool,
		router:      router,
		policyChain: policyChain,
	}, nil
}

// Start starts the application
func (a *App) Start() error {
	a.logger.Info("Starting server",
		zap.String("host", a.config.Server.Host),
		zap.Int("port", a.config.Server.Port))
	return a.server.ListenAndServe()
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() error {
	a.logger.Info("Shutting down server")
	return a.server.Close()
} 