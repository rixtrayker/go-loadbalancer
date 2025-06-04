package middleware

import (
	"net/http"
	"time"

	"github.com/rixtrayker/go-loadbalancer/internal/monitoring"
	"github.com/rixtrayker/go-loadbalancer/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// MonitoringMiddleware wraps an http.Handler with monitoring capabilities
func MonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start span
		ctx, span := tracing.StartSpan(r.Context(), "http.request",
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.user_agent", r.UserAgent()),
			),
		)
		defer span.End()

		// Create response writer wrapper for status code
		rw := newResponseWriter(w)

		// Start timer for request duration
		start := time.Now()

		// Process request
		next.ServeHTTP(rw, r.WithContext(ctx))

		// Record metrics
		duration := time.Since(start).Seconds()
		monitoring.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		monitoring.RequestTotal.WithLabelValues(r.Method, r.URL.Path, string(rw.statusCode)).Inc()

		// Record span attributes
		span.SetAttributes(
			attribute.Int("http.status_code", rw.statusCode),
			attribute.Float64("http.duration_seconds", duration),
		)

		// Record error if status code indicates error
		if rw.statusCode >= 400 {
			span.RecordError(nil, trace.WithAttributes(
				attribute.Int("http.status_code", rw.statusCode),
			))
		}
	})
}

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

// ErrorMiddleware wraps an http.Handler with error handling capabilities
func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Record error in span
				span := tracing.SpanFromContext(r.Context())
				span.RecordError(nil, trace.WithAttributes(
					attribute.String("error.type", "panic"),
					attribute.String("error.message", err.(string)),
				))

				// Record metric
				monitoring.RequestTotal.WithLabelValues(r.Method, r.URL.Path, "500").Inc()

				// Return 500 Internal Server Error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
} 