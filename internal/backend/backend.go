package backend

import (
	"net/url"
	"sync"
	"sync/atomic"
)

// Backend represents a backend server
type Backend struct {
	URL           *url.URL
	Weight        int
	Healthy       bool
	ActiveConns   int32
	TotalRequests int64
	mutex         sync.RWMutex
}

// NewBackend creates a new backend instance
func NewBackend(urlStr string, weight int) (*Backend, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &Backend{
		URL:     url,
		Weight:  weight,
		Healthy: true,
	}, nil
}

// IsHealthy returns the health status of the backend
func (b *Backend) IsHealthy() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.Healthy
}

// SetHealth sets the health status of the backend
func (b *Backend) SetHealth(healthy bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.Healthy = healthy
}

// IncrementConnections increments the active connection count
func (b *Backend) IncrementConnections() {
	atomic.AddInt32(&b.ActiveConns, 1)
}

// DecrementConnections decrements the active connection count
func (b *Backend) DecrementConnections() {
	atomic.AddInt32(&b.ActiveConns, -1)
}

// GetActiveConnections returns the number of active connections
func (b *Backend) GetActiveConnections() int {
	return int(atomic.LoadInt32(&b.ActiveConns))
}

// IncrementRequests increments the total request count
func (b *Backend) IncrementRequests() {
	atomic.AddInt64(&b.TotalRequests, 1)
}

// GetTotalRequests returns the total number of requests
func (b *Backend) GetTotalRequests() int64 {
	return atomic.LoadInt64(&b.TotalRequests)
}
