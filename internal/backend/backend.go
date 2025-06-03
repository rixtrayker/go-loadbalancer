package backend

import (
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// State represents the current state of a backend server
type State int

const (
	// StateUnknown indicates the backend state is not yet determined
	StateUnknown State = iota
	// StateHealthy indicates the backend is healthy and can receive traffic
	StateHealthy
	// StateUnhealthy indicates the backend is unhealthy and should not receive traffic
	StateUnhealthy
)

// Backend represents a backend server with its state and metrics
type Backend struct {
	URL           *url.URL
	Name          string
	Weight        int
	MaxConns      int
	CurrentConns  int32
	State         State
	LastChecked   time.Time
	ResponseTime  time.Duration
	FailCount     int32
	SuccessCount  int32
	mu            sync.RWMutex
	Active        bool
	Healthy       bool
	LastCheck     time.Time
	CircuitBreaker *CircuitBreaker
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	Failures       int
	Threshold      int
	ResetTimeout   time.Duration
	LastFailure    time.Time
	State          CircuitState
	mu             sync.RWMutex
}

type CircuitState int

const (
	Closed CircuitState = iota
	Open
	HalfOpen
)

// NewBackend creates a new backend instance
func NewBackend(name string, urlStr string, weight, maxConns int) (*Backend, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &Backend{
		URL:          parsedURL,
		Name:         name,
		Weight:       weight,
		MaxConns:     maxConns,
		State:        StateUnknown,
		LastChecked:  time.Now(),
		ResponseTime: 0,
		Active:       true,
		Healthy:      true,
		LastCheck:    time.Now(),
		CircuitBreaker: &CircuitBreaker{
			Failures:     0,
			Threshold:    5,
			ResetTimeout: 30 * time.Second,
			State:        Closed,
		},
	}, nil
}

// IncrementConns increments the current connection count
func (b *Backend) IncrementConns() bool {
	current := atomic.LoadInt32(&b.CurrentConns)
	if current >= int32(b.MaxConns) {
		return false
	}
	atomic.AddInt32(&b.CurrentConns, 1)
	return true
}

// DecrementConns decrements the current connection count
func (b *Backend) DecrementConns() {
	atomic.AddInt32(&b.CurrentConns, -1)
}

// UpdateState updates the backend state
func (b *Backend) UpdateState(state State) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.State = state
	b.LastChecked = time.Now()
}

// UpdateResponseTime updates the response time metric
func (b *Backend) UpdateResponseTime(duration time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ResponseTime = duration
}

// IncrementFailCount increments the failure counter
func (b *Backend) IncrementFailCount() {
	atomic.AddInt32(&b.FailCount, 1)
}

// IncrementSuccessCount increments the success counter
func (b *Backend) IncrementSuccessCount() {
	atomic.AddInt32(&b.SuccessCount, 1)
}

// IsAvailable checks if the backend is available for new connections
func (b *Backend) IsAvailable() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Active && b.Healthy && b.CurrentConns < int32(b.MaxConns) && 
		b.CircuitBreaker.State != Open
}

// IncrementConn increases the current connection count
func (b *Backend) IncrementConn() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.CurrentConns++
}

// DecrementConn decreases the current connection count
func (b *Backend) DecrementConn() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.CurrentConns > 0 {
		b.CurrentConns--
	}
}

// RecordFailure records a failure and updates circuit breaker state
func (b *Backend) RecordFailure() {
	b.CircuitBreaker.mu.Lock()
	defer b.CircuitBreaker.mu.Unlock()

	b.CircuitBreaker.Failures++
	b.CircuitBreaker.LastFailure = time.Now()

	if b.CircuitBreaker.Failures >= b.CircuitBreaker.Threshold {
		b.CircuitBreaker.State = Open
	}
}

// RecordSuccess resets the circuit breaker on success
func (b *Backend) RecordSuccess() {
	b.CircuitBreaker.mu.Lock()
	defer b.CircuitBreaker.mu.Unlock()

	b.CircuitBreaker.Failures = 0
	b.CircuitBreaker.State = Closed
}

// CheckCircuitBreaker checks and potentially updates the circuit breaker state
func (b *Backend) CheckCircuitBreaker() {
	b.CircuitBreaker.mu.Lock()
	defer b.CircuitBreaker.mu.Unlock()

	if b.CircuitBreaker.State == Open {
		if time.Since(b.CircuitBreaker.LastFailure) > b.CircuitBreaker.ResetTimeout {
			b.CircuitBreaker.State = HalfOpen
		}
	}
} 