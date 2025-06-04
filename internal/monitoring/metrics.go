package monitoring

import (
	"context"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rixtrayker/go-loadbalancer/internal/logging"
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
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	// Backend metrics
	BackendHealthStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "loadbalancer_backend_health_status",
			Help: "Health status of backend servers (1 = healthy, 0 = unhealthy)",
		},
		[]string{"backend", "pool"},
	)

	BackendResponseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loadbalancer_backend_response_time_seconds",
			Help:    "Backend response time in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"backend", "pool"},
	)

	BackendErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_backend_errors_total",
			Help: "Total number of backend errors",
		},
		[]string{"backend", "pool", "error_type"},
	)

	// Policy metrics
	PolicyViolations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_policy_violations_total",
			Help: "Total number of policy violations",
		},
		[]string{"policy_type", "path"},
	)

	RateLimitHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_rate_limit_hits_total",
			Help: "Total number of rate limit hits",
		},
		[]string{"policy", "path"},
	)

	// Connection metrics
	ActiveConnections = promauto.NewGaugeVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_active_connections",
			Help: "Number of active connections",
		},
		[]string{"backend", "pool"},
	)

	ConnectionErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_connection_errors_total",
			Help: "Total number of connection errors",
		},
		[]string{"backend", "pool", "error_type"},
	)

	// System metrics
	MemoryUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "loadbalancer_memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
		[]string{"type"},
	)

	CPUUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "loadbalancer_cpu_usage_percent",
			Help: "Current CPU usage percentage",
		},
	)

	GoroutineCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "loadbalancer_goroutine_count",
			Help: "Current number of goroutines",
		},
	)

	// Cache metrics
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache"},
	)

	CacheSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "loadbalancer_cache_size_bytes",
			Help: "Current cache size in bytes",
		},
		[]string{"cache"},
	)
)

// MetricsCollector collects system metrics
type MetricsCollector struct {
	logger *logging.Logger
	done   chan struct{}
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(logger *logging.Logger) *MetricsCollector {
	return &MetricsCollector{
		logger: logger,
		done:   make(chan struct{}),
	}
}

// Start begins collecting system metrics
func (mc *MetricsCollector) Start(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				mc.collectSystemMetrics()
			case <-ctx.Done():
				mc.logger.Info("Stopping metrics collector")
				close(mc.done)
				return
			}
		}
	}()
}

// Stop stops the metrics collector
func (mc *MetricsCollector) Stop() {
	<-mc.done
}

// collectSystemMetrics collects system metrics
func (mc *MetricsCollector) collectSystemMetrics() {
	// Collect memory stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	MemoryUsage.WithLabelValues("heap_alloc").Set(float64(memStats.HeapAlloc))
	MemoryUsage.WithLabelValues("heap_sys").Set(float64(memStats.HeapSys))
	MemoryUsage.WithLabelValues("heap_idle").Set(float64(memStats.HeapIdle))
	MemoryUsage.WithLabelValues("heap_inuse").Set(float64(memStats.HeapInuse))
	MemoryUsage.WithLabelValues("stack_inuse").Set(float64(memStats.StackInuse))
	MemoryUsage.WithLabelValues("stack_sys").Set(float64(memStats.StackSys))

	// Collect goroutine count
	GoroutineCount.Set(float64(runtime.NumGoroutine()))

	// Note: CPU usage requires additional libraries like gopsutil
	// This is a simplified version
	CPUUsage.Set(float64(runtime.NumCPU()))
}

// RecordRequestMetrics records request metrics
func RecordRequestMetrics(method, path, status string, duration time.Duration) {
	RequestTotal.WithLabelValues(method, path, status).Inc()
	RequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// RecordBackendMetrics records backend metrics
func RecordBackendMetrics(backend, pool string, healthy bool, responseTime time.Duration) {
	var healthStatus float64
	if healthy {
		healthStatus = 1
	}
	BackendHealthStatus.WithLabelValues(backend, pool).Set(healthStatus)
	BackendResponseTime.WithLabelValues(backend, pool).Observe(responseTime.Seconds())
}

// RecordBackendError records a backend error
func RecordBackendError(backend, pool, errorType string) {
	BackendErrors.WithLabelValues(backend, pool, errorType).Inc()
}

// RecordPolicyViolation records a policy violation
func RecordPolicyViolation(policyType, path string) {
	PolicyViolations.WithLabelValues(policyType, path).Inc()
}

// RecordRateLimitHit records a rate limit hit
func RecordRateLimitHit(policy, path string) {
	RateLimitHits.WithLabelValues(policy, path).Inc()
}

// RecordCacheMetrics records cache metrics
func RecordCacheMetrics(cache string, hit bool, size int64) {
	if hit {
		CacheHits.WithLabelValues(cache).Inc()
	} else {
		CacheMisses.WithLabelValues(cache).Inc()
	}
	CacheSize.WithLabelValues(cache).Set(float64(size))
}
