package healthcheck

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/amr/go-loadbalancer/internal/backend"
)

// HealthChecker manages health checks for backends
type HealthChecker struct {
	backends    []*backend.Backend
	interval    time.Duration
	timeout     time.Duration
	stopChan    chan struct{}
	probeType   string
	mu          sync.RWMutex
}

// NewHealthChecker creates a new health checker instance
func NewHealthChecker(interval, timeout time.Duration, probeType string) *HealthChecker {
	return &HealthChecker{
		interval:  interval,
		timeout:   timeout,
		stopChan:  make(chan struct{}),
		probeType: probeType,
	}
}

// AddBackend adds a backend to be health checked
func (hc *HealthChecker) AddBackend(b *backend.Backend) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.backends = append(hc.backends, b)
}

// RemoveBackend removes a backend from health checking
func (hc *HealthChecker) RemoveBackend(url string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	
	for i, b := range hc.backends {
		if b.URL == url {
			hc.backends = append(hc.backends[:i], hc.backends[i+1:]...)
			return
		}
	}
}

// Start begins the health checking process
func (hc *HealthChecker) Start() {
	ticker := time.NewTicker(hc.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				hc.checkAllBackends()
			case <-hc.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the health checking process
func (hc *HealthChecker) Stop() {
	close(hc.stopChan)
}

// checkAllBackends performs health checks on all registered backends
func (hc *HealthChecker) checkAllBackends() {
	hc.mu.RLock()
	backends := make([]*backend.Backend, len(hc.backends))
	copy(backends, hc.backends)
	hc.mu.RUnlock()

	for _, b := range backends {
		go hc.checkBackend(b)
	}
}

// checkBackend performs a health check on a single backend
func (hc *HealthChecker) checkBackend(b *backend.Backend) {
	ctx, cancel := context.WithTimeout(context.Background(), hc.timeout)
	defer cancel()

	var healthy bool
	var err error

	switch hc.probeType {
	case "http":
		healthy, err = checkHTTP(ctx, b.URL)
	case "tcp":
		healthy, err = checkTCP(ctx, b.URL)
	default:
		healthy, err = checkHTTP(ctx, b.URL)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if err != nil || !healthy {
		b.Healthy = false
		b.RecordFailure()
	} else {
		b.Healthy = true
		b.RecordSuccess()
	}
	b.LastCheck = time.Now()
}

// checkHTTP performs an HTTP health check
func checkHTTP(ctx context.Context, url string) (bool, error) {
	// TODO: Implement actual HTTP health check logic
	// This could involve making an HTTP request and checking the response
	// For now, we'll just log that we're checking
	log.Printf("Checking health of backend %s via HTTP", url)
	return true, nil
}

// checkTCP performs a TCP health check
func checkTCP(ctx context.Context, url string) (bool, error) {
	// TODO: Implement actual TCP health check logic
	// This could involve establishing a TCP connection and checking if it's open
	// For now, we'll just log that we're checking
	log.Printf("Checking health of backend %s via TCP", url)
	return true, nil
} 