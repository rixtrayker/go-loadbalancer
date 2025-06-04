package tracer

import (
	"io"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

// InitTracer initializes the Jaeger tracer
func InitTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	// Check if tracing is enabled
	if os.Getenv("TRACING_ENABLED") != "true" {
		return opentracing.NoopTracer{}, &noopCloser{}, nil
	}

	// Configure Jaeger
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: os.Getenv("JAEGER_AGENT_HOST") + ":" + os.Getenv("JAEGER_AGENT_PORT"),
		},
	}

	// Initialize tracer
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		return nil, nil, err
	}

	// Set as global tracer
	opentracing.SetGlobalTracer(tracer)

	return tracer, closer, nil
}

// noopCloser is a no-op io.Closer
type noopCloser struct{}

func (n *noopCloser) Close() error {
	return nil
}
