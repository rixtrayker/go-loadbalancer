package middleware

import (
	"net/http"
	"time"

	"github.com/rixtrayker/go-loadbalancer/internal/logging"
	"github.com/rixtrayker/go-loadbalancer/internal/monitoring"
)

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// MonitoringMiddleware wraps an http.Handler with monitoring capabilities
func MonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create response writer wrapper for status code
		rw := newResponseWriter(w)

		// Start timer for request duration
		start := time.Now()

		// Process request
		next.ServeHTTP(rw, r)

		// Record metrics
		duration := time.Since(start)
		monitoring.RecordRequestMetrics(r.Method, r.URL.Path, string(rw.statusCode), duration)
	})
}

// LoggingMiddleware wraps an http.Handler with logging capabilities
func LoggingMiddleware(logger *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create response writer wrapper for status code
			rw := newResponseWriter(w)

			// Start timer for request duration
			start := time.Now()

			// Log request
			logger.Info("Request started",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)

			// Process request
			next.ServeHTTP(rw, r)

			// Log response
			logger.Info("Request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration", time.Since(start).String(),
			)
		})
	}
}

// TracingMiddleware wraps an http.Handler with tracing capabilities
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simplified implementation without OpenTelemetry
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware wraps an http.Handler with panic recovery
func RecoveryMiddleware(logger *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log error
					logger.Error("Panic recovered",
						"error", err,
						"method", r.Method,
						"path", r.URL.Path,
					)

					// Return 500 Internal Server Error
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
