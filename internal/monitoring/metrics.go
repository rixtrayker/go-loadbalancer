package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Request metrics
	RequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_requests_total",
			Help: "Total number of requests processed",
		},
		[]string{"method", "path", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loadbalancer_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Backend metrics
	BackendHealthStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "loadbalancer_backend_health_status",
			Help: "Health status of backend servers (1 = healthy, 0 = unhealthy)",
		},
		[]string{"backend"},
	)

	BackendResponseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loadbalancer_backend_response_time_seconds",
			Help:    "Backend response time in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"backend"},
	)

	BackendErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_backend_errors_total",
			Help: "Total number of backend errors",
		},
		[]string{"backend", "error_type"},
	)

	// Policy metrics
	PolicyViolations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_policy_violations_total",
			Help: "Total number of policy violations",
		},
		[]string{"policy_type"},
	)

	RateLimitHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_rate_limit_hits_total",
			Help: "Total number of rate limit hits",
		},
		[]string{"policy"},
	)

	// Connection metrics
	ActiveConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "loadbalancer_active_connections",
			Help: "Number of active connections",
		},
		[]string{"backend"},
	)

	ConnectionErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_connection_errors_total",
			Help: "Total number of connection errors",
		},
		[]string{"backend", "error_type"},
	)

	// System metrics
	MemoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "loadbalancer_memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
	)

	CPUUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "loadbalancer_cpu_usage_percent",
			Help: "Current CPU usage percentage",
		},
	)

	// Cache metrics
	CacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "loadbalancer_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "loadbalancer_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	CacheSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "loadbalancer_cache_size_bytes",
			Help: "Current cache size in bytes",
		},
	)
) 