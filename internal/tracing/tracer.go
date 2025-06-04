package tracing

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rixtrayker/go-loadbalancer/configs"
	"github.com/rixtrayker/go-loadbalancer/internal/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DefaultConfig returns the default tracing configuration
func DefaultConfig() configs.TracingConfig {
	return configs.TracingConfig{
		Enabled:        false,
		ServiceName:    "go-loadbalancer",
		ServiceVersion: "0.1.0",
		Environment:    "development",
		Endpoint:       "localhost:4317",
		SamplingRate:   1.0,
		Protocol:       "grpc",
		Secure:         false,
	}
}

// applyDefaults applies default values to the configuration if they are not set
func applyDefaults(config *configs.TracingConfig) {
	defaultConfig := DefaultConfig()
	if config.ServiceName == "" {
		config.ServiceName = defaultConfig.ServiceName
	}
	if config.ServiceVersion == "" {
		config.ServiceVersion = defaultConfig.ServiceVersion
	}
	if config.Environment == "" {
		config.Environment = defaultConfig.Environment
	}
	if config.Endpoint == "" {
		config.Endpoint = defaultConfig.Endpoint
	}
	if config.SamplingRate == 0 {
		config.SamplingRate = defaultConfig.SamplingRate
	}
	if config.Protocol == "" {
		config.Protocol = defaultConfig.Protocol
	}
}

// Tracer represents the OpenTelemetry tracer
type Tracer struct {
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
	logger         *logging.Logger
}

// NewTracer creates a new tracer instance
func NewTracer(config configs.TracingConfig, logger *logging.Logger) (*Tracer, error) {
	if !config.Enabled {
		return &Tracer{
			tracer: otel.Tracer("noop"),
			logger: logger,
		}, nil
	}

	// Apply default values
	applyDefaults(&config)

	// Create resource with service information
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			attribute.String("environment", config.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create exporter based on protocol
	var exporter sdktrace.SpanExporter
	if config.Protocol == "http" {
		// HTTP exporter
		httpOpts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(config.Endpoint),
		}
		
		if !config.Secure {
			httpOpts = append(httpOpts, otlptracehttp.WithInsecure())
		}
		
		httpExporter, err := otlptrace.New(
			context.Background(),
			otlptracehttp.NewClient(httpOpts...),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP HTTP exporter: %w", err)
		}
		exporter = httpExporter
	} else {
		// Default to gRPC exporter
		grpcOpts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(config.Endpoint),
		}
		
		if !config.Secure {
			grpcOpts = append(grpcOpts, otlptracegrpc.WithInsecure())
			grpcOpts = append(grpcOpts, otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
		}
		
		grpcExporter, err := otlptrace.New(
			context.Background(),
			otlptracegrpc.NewClient(grpcOpts...),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP gRPC exporter: %w", err)
		}
		exporter = grpcExporter
	}

	// Create batch span processor
	bsp := sdktrace.NewBatchSpanProcessor(exporter)

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(config.SamplingRate)),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(res),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)
	
	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Tracer{
		tracer:         tp.Tracer(config.ServiceName),
		tracerProvider: tp,
		logger:         logger,
	}, nil
}

// StartSpan starts a new span
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// Shutdown gracefully shuts down the tracer
func (t *Tracer) Shutdown(ctx context.Context) error {
	if t.tracerProvider == nil {
		return nil
	}
	
	return t.tracerProvider.Shutdown(ctx)
}

// TracingMiddleware wraps an http.Handler with tracing capabilities
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract trace context from request
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		// Start span
		ctx, span := otel.Tracer("http.server").Start(ctx, r.URL.Path,
			trace.WithAttributes(
				semconv.HTTPMethod(r.Method),
				semconv.HTTPTarget(r.URL.Path),
				semconv.HTTPRoute(r.URL.Path),
				semconv.HTTPUserAgent(r.UserAgent()),
				semconv.HTTPClientIP(r.RemoteAddr),
				attribute.String("http.host", r.Host),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Create response writer wrapper for status code
		rw := newResponseWriter(w)

		// Start timer for request duration
		start := time.Now()

		// Process request
		next.ServeHTTP(rw, r.WithContext(ctx))

		// Record span attributes
		span.SetAttributes(
			semconv.HTTPStatusCode(rw.statusCode),
			attribute.Float64("http.duration_seconds", time.Since(start).Seconds()),
		)

		// Record error if status code indicates error
		if rw.statusCode >= 400 {
			span.SetStatus(codes.Error, fmt.Sprintf("%s: %s", rw.statusCode, http.StatusText(rw.statusCode)))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	})
}

// SpanFromContext returns the current span from context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetSpanAttributes sets attributes on the current span
func SetSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// RecordError records an error on the current span
func RecordError(ctx context.Context, err error, opts ...trace.EventOption) {
	span := SpanFromContext(ctx)
	span.RecordError(err, opts...)
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
