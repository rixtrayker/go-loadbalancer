package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics collects and exposes application metrics
type Metrics struct {
	requestCounter   *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	backendLatency   *prometheus.HistogramVec
	activeConnections *prometheus.GaugeVec
	backendHealth    *prometheus.GaugeVec
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	m := &Metrics{
		requestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "loadbalancer_requests_total",
				Help: "Total number of requests processed",
			},
			[]string{"path", "method", "status"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "loadbalancer_request_duration_seconds",
				Help:    "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"path", "method"},
		),
		backendLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "loadbalancer_backend_latency_seconds",
				Help:    "Backend request latency in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"backend"},
		),
		activeConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "loadbalancer_active_connections",
				Help: "Number of active connections",
			},
			[]string{"backend"},
		),
		backendHealth: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "loadbalancer_backend_health",
				Help: "Backend health status (1 = healthy, 0 = unhealthy)",
			},
			[]string{"backend"},
		),
	}

	return m
}

// RecordRequest records a processed request
func (m *Metrics) RecordRequest(path, method, status string) {
	m.requestCounter.WithLabelValues(path, method, status).Inc()
}

// RecordRequestDuration records the duration of a request
func (m *Metrics) RecordRequestDuration(path, method string, duration time.Duration) {
	m.requestDuration.WithLabelValues(path, method).Observe(duration.Seconds())
}

// RecordBackendLatency records the latency of a backend request
func (m *Metrics) RecordBackendLatency(backend string, duration time.Duration) {
	m.backendLatency.WithLabelValues(backend).Observe(duration.Seconds())
}

// SetActiveConnections sets the number of active connections for a backend
func (m *Metrics) SetActiveConnections(backend string, count int) {
	m.activeConnections.WithLabelValues(backend).Set(float64(count))
}

// SetBackendHealth sets the health status of a backend
func (m *Metrics) SetBackendHealth(backend string, healthy bool) {
	var value float64
	if healthy {
		value = 1
	}
	m.backendHealth.WithLabelValues(backend).Set(value)
}
