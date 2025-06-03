package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// MetricType represents the type of metric
type MetricType int

const (
	// CounterType represents a counter metric
	CounterType MetricType = iota
	// GaugeType represents a gauge metric
	GaugeType
	// HistogramType represents a histogram metric
	HistogramType
)

// Metric represents a metric
type Metric struct {
	Name   string
	Type   MetricType
	Value  float64
	Labels map[string]string
}

// Counter represents a counter metric
type Counter struct {
	name   string
	value  uint64
	labels map[string]string
}

// Gauge represents a gauge metric
type Gauge struct {
	name   string
	value  float64
	labels map[string]string
	mu     sync.RWMutex
}

// Histogram represents a histogram metric
type Histogram struct {
	name    string
	buckets []float64
	values  []uint64
	labels  map[string]string
	mu      sync.RWMutex
}

// Registry represents a metrics registry
type Registry struct {
	counters   map[string]*Counter
	gauges     map[string]*Gauge
	histograms map[string]*Histogram
	mu         sync.RWMutex
}

// Metrics represents load balancer metrics
type Metrics struct {
	mu sync.RWMutex

	// Request metrics
	TotalRequests     uint64
	SuccessfulRequests uint64
	FailedRequests    uint64
	RequestLatency    time.Duration

	// Backend metrics
	ActiveBackends    int32
	HealthyBackends   int32
	TotalConnections  int32
	BackendLatencies  map[string]time.Duration
}

// NewRegistry creates a new metrics registry
func NewRegistry() *Registry {
	return &Registry{
		counters:   make(map[string]*Counter),
		gauges:     make(map[string]*Gauge),
		histograms: make(map[string]*Histogram),
	}
}

// NewCounter creates a new counter metric
func (r *Registry) NewCounter(name string, labels map[string]string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()

	counter := &Counter{
		name:   name,
		labels: labels,
	}
	r.counters[name] = counter
	return counter
}

// NewGauge creates a new gauge metric
func (r *Registry) NewGauge(name string, labels map[string]string) *Gauge {
	r.mu.Lock()
	defer r.mu.Unlock()

	gauge := &Gauge{
		name:   name,
		labels: labels,
	}
	r.gauges[name] = gauge
	return gauge
}

// NewHistogram creates a new histogram metric
func (r *Registry) NewHistogram(name string, buckets []float64, labels map[string]string) *Histogram {
	r.mu.Lock()
	defer r.mu.Unlock()

	histogram := &Histogram{
		name:    name,
		buckets: buckets,
		values:  make([]uint64, len(buckets)),
		labels:  labels,
	}
	r.histograms[name] = histogram
	return histogram
}

// Increment increments a counter
func (c *Counter) Increment() {
	atomic.AddUint64(&c.value, 1)
}

// Add adds a value to a counter
func (c *Counter) Add(value uint64) {
	atomic.AddUint64(&c.value, value)
}

// Get returns the counter value
func (c *Counter) Get() uint64 {
	return atomic.LoadUint64(&c.value)
}

// Set sets a gauge value
func (g *Gauge) Set(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value = value
}

// Get returns the gauge value
func (g *Gauge) Get() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.value
}

// Observe adds a value to a histogram
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i, bucket := range h.buckets {
		if value <= bucket {
			h.values[i]++
		}
	}
}

// GetBuckets returns the histogram buckets and their values
func (h *Histogram) GetBuckets() ([]float64, []uint64) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.buckets, h.values
}

// Collect returns all metrics
func (r *Registry) Collect() []Metric {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metrics := make([]Metric, 0)

	// Collect counters
	for _, counter := range r.counters {
		metrics = append(metrics, Metric{
			Name:   counter.name,
			Type:   CounterType,
			Value:  float64(counter.Get()),
			Labels: counter.labels,
		})
	}

	// Collect gauges
	for _, gauge := range r.gauges {
		metrics = append(metrics, Metric{
			Name:   gauge.name,
			Type:   GaugeType,
			Value:  gauge.Get(),
			Labels: gauge.labels,
		})
	}

	// Collect histograms
	for _, histogram := range r.histograms {
		buckets, values := histogram.GetBuckets()
		for i := range buckets {
			metrics = append(metrics, Metric{
				Name:   histogram.name + "_bucket",
				Type:   HistogramType,
				Value:  float64(values[i]),
				Labels: histogram.labels,
			})
		}
	}

	return metrics
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		BackendLatencies: make(map[string]time.Duration),
	}
}

// RecordRequest records a request
func (m *Metrics) RecordRequest(success bool, latency time.Duration) {
	atomic.AddUint64(&m.TotalRequests, 1)
	if success {
		atomic.AddUint64(&m.SuccessfulRequests, 1)
	} else {
		atomic.AddUint64(&m.FailedRequests, 1)
	}

	m.mu.Lock()
	m.RequestLatency = (m.RequestLatency + latency) / 2
	m.mu.Unlock()
}

// UpdateBackendMetrics updates backend metrics
func (m *Metrics) UpdateBackendMetrics(active, healthy, connections int32) {
	atomic.StoreInt32(&m.ActiveBackends, active)
	atomic.StoreInt32(&m.HealthyBackends, healthy)
	atomic.StoreInt32(&m.TotalConnections, connections)
}

// RecordBackendLatency records latency for a backend
func (m *Metrics) RecordBackendLatency(backend string, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if current, exists := m.BackendLatencies[backend]; exists {
		m.BackendLatencies[backend] = (current + latency) / 2
	} else {
		m.BackendLatencies[backend] = latency
	}
}

// GetMetrics returns the current metrics
func (m *Metrics) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"total_requests":      atomic.LoadUint64(&m.TotalRequests),
		"successful_requests": atomic.LoadUint64(&m.SuccessfulRequests),
		"failed_requests":     atomic.LoadUint64(&m.FailedRequests),
		"request_latency":     m.RequestLatency,
		"active_backends":     atomic.LoadInt32(&m.ActiveBackends),
		"healthy_backends":    atomic.LoadInt32(&m.HealthyBackends),
		"total_connections":   atomic.LoadInt32(&m.TotalConnections),
		"backend_latencies":   m.BackendLatencies,
	}
} 