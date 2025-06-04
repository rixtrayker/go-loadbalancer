package ratelimit

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RateLimiter implements rate limiting
type RateLimiter struct {
	limits     map[string]*limit
	mutex      sync.RWMutex
	cleanupInt time.Duration
}

type limit struct {
	tokens     int
	lastRefill time.Time
	rate       int
	per        time.Duration
}

var (
	// Global rate limiter instance
	globalLimiter = NewRateLimiter()
)

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		limits:     make(map[string]*limit),
		cleanupInt: 10 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Apply applies rate limiting to a request
func Apply(rateStr string, r *http.Request) error {
	// Parse rate limit string (e.g., "100/minute", "10/second")
	parts := strings.Split(rateStr, "/")
	if len(parts) != 2 {
		return errors.New("invalid rate limit format")
	}

	rate, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid rate: %w", err)
	}

	var per time.Duration
	switch strings.ToLower(parts[1]) {
	case "second":
		per = time.Second
	case "minute":
		per = time.Minute
	case "hour":
		per = time.Hour
	default:
		return errors.New("invalid time unit")
	}

	// Get client IP as key
	key := getClientIP(r)

	// Check rate limit
	return globalLimiter.Allow(key, rate, per)
}

// Allow checks if a request is allowed based on rate limits
func (rl *RateLimiter) Allow(key string, rate int, per time.Duration) error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	l, exists := rl.limits[key]
	if !exists {
		// Create new limit for this key
		l = &limit{
			tokens:     rate,
			lastRefill: time.Now(),
			rate:       rate,
			per:        per,
		}
		rl.limits[key] = l
	} else {
		// Refill tokens based on elapsed time
		now := time.Now()
		elapsed := now.Sub(l.lastRefill)
		tokensToAdd := int(float64(elapsed) / float64(per) * float64(rate))
		
		if tokensToAdd > 0 {
			l.tokens = min(l.rate, l.tokens+tokensToAdd)
			l.lastRefill = now
		}
	}

	// Check if we have tokens available
	if l.tokens <= 0 {
		return errors.New("rate limit exceeded")
	}

	// Consume a token
	l.tokens--
	return nil
}

// cleanup periodically removes expired limits
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInt)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for key, l := range rl.limits {
			// Remove limits that haven't been used for a while
			if now.Sub(l.lastRefill) > l.per*2 {
				delete(rl.limits, key)
			}
		}
		rl.mutex.Unlock()
	}
}

// getClientIP extracts the client IP from a request
func getClientIP(r *http.Request) string {
	// Try X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Try X-Real-IP header
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
