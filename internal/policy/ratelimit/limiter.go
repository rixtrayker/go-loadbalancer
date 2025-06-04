package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/amr/go-loadbalancer/internal/policy"
)

// RateLimiter implements rate limiting policy
type RateLimiter struct {
	policy.BasePolicy
	rate       int
	window     time.Duration
	requests   map[string][]time.Time
	mu         sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		BasePolicy: policy.BasePolicy{
			PolicyName: "rate-limiter",
		},
		rate:     rate,
		window:   window,
		requests: make(map[string][]time.Time),
	}
}

// Apply implements the Policy interface
func (rl *RateLimiter) Apply(req *http.Request, _ *http.Response) error {
	key := rl.getKey(req)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Clean up old requests
	if requests, exists := rl.requests[key]; exists {
		var valid []time.Time
		for _, t := range requests {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		rl.requests[key] = valid
	}

	// Check if rate limit is exceeded
	if len(rl.requests[key]) >= rl.rate {
		return ErrRateLimitExceeded
	}

	// Add current request
	rl.requests[key] = append(rl.requests[key], now)
	return nil
}

// getKey generates a key for rate limiting
func (rl *RateLimiter) getKey(req *http.Request) string {
	// Use IP address as the key
	return req.RemoteAddr
}

// ErrRateLimitExceeded is returned when rate limit is exceeded
var ErrRateLimitExceeded = &RateLimitError{}

// RateLimitError represents a rate limit error
type RateLimitError struct{}

func (e *RateLimitError) Error() string {
	return "rate limit exceeded"
}

// Limiter implements a rate limiter using the token bucket algorithm
type Limiter struct {
	rate       float64 // tokens per second
	bucketSize float64
	tokens     float64
	lastUpdate time.Time
	mu         sync.Mutex
}

// NewLimiter creates a new token bucket limiter
func NewLimiter(rate, bucketSize float64) *Limiter {
	return &Limiter{
		rate:       rate,
		bucketSize: bucketSize,
		tokens:     bucketSize,
		lastUpdate: time.Now(),
	}
}

// Allow checks if a request is allowed under the rate limit
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastUpdate).Seconds()
	l.lastUpdate = now

	// Add new tokens based on elapsed time
	l.tokens = min(l.bucketSize, l.tokens+elapsed*l.rate)

	// Check if we have enough tokens
	if l.tokens >= 1.0 {
		l.tokens -= 1.0
		return true
	}

	return false
}

// min returns the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// TokenBucketLimiter represents a rate limiter for multiple keys using token bucket algorithm
type TokenBucketLimiter struct {
	limiters map[string]*Limiter
	mu       sync.RWMutex
	rate     float64
	burst    int
}

// NewTokenBucketLimiter creates a new token bucket rate limiter for multiple keys
func NewTokenBucketLimiter(rate float64, burst int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		limiters: make(map[string]*Limiter),
		rate:     rate,
		burst:    burst,
	}
}

// Allow checks if a request is allowed for the given key
func (rl *TokenBucketLimiter) Allow(key string) bool {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter = NewLimiter(rl.rate, float64(rl.burst))
		rl.limiters[key] = limiter
		rl.mu.Unlock()
	}

	return limiter.Allow()
}

// Remove removes a rate limiter for a specific key
func (rl *TokenBucketLimiter) Remove(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.limiters, key)
}
