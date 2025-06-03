package tracer

import (
	"context"
	"time"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

// Span represents a tracing span
type Span struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Tags      map[string]string
	Events    []Event
	jaegerSpan jaeger.Span
}

// Event represents a span event
type Event struct {
	Name      string
	Timestamp time.Time
	Tags      map[string]string
}

// Tracer represents a distributed tracer
type Tracer interface {
	// StartSpan starts a new span
	StartSpan(ctx context.Context, name string) (context.Context, Span)
	// EndSpan ends a span
	EndSpan(span Span)
	// AddEvent adds an event to a span
	AddEvent(span *Span, name string, tags map[string]string)
	// AddTag adds a tag to a span
	AddTag(span *Span, key, value string)
}

// JaegerTracer implements the Tracer interface using Jaeger
type JaegerTracer struct {
	tracer jaeger.Tracer
}

// NewJaegerTracer creates a new Jaeger tracer
func NewJaegerTracer(serviceName string) (*JaegerTracer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, _, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
		config.Metrics(metrics.NullFactory),
	)
	if err != nil {
		return nil, err
	}

	return &JaegerTracer{
		tracer: tracer,
	}, nil
}

// StartSpan implements the Tracer interface
func (t *JaegerTracer) StartSpan(ctx context.Context, name string) (context.Context, Span) {
	jaegerSpan := t.tracer.StartSpan(name)
	ctx = jaeger.ContextWithSpan(ctx, jaegerSpan)

	return ctx, Span{
		Name:      name,
		StartTime: time.Now(),
		Tags:      make(map[string]string),
		Events:    make([]Event, 0),
		jaegerSpan: jaegerSpan,
	}
}

// EndSpan implements the Tracer interface
func (t *JaegerTracer) EndSpan(span Span) {
	span.EndTime = time.Now()
	span.jaegerSpan.Finish()
}

// AddEvent implements the Tracer interface
func (t *JaegerTracer) AddEvent(span *Span, name string, tags map[string]string) {
	span.Events = append(span.Events, Event{
		Name:      name,
		Timestamp: time.Now(),
		Tags:      tags,
	})
	span.jaegerSpan.LogKV("event", name, "tags", tags)
}

// AddTag implements the Tracer interface
func (t *JaegerTracer) AddTag(span *Span, key, value string) {
	if span.Tags == nil {
		span.Tags = make(map[string]string)
	}
	span.Tags[key] = value
	span.jaegerSpan.SetTag(key, value)
}

// NoopTracer implements a no-op tracer
type NoopTracer struct{}

// StartSpan implements the Tracer interface
func (t *NoopTracer) StartSpan(ctx context.Context, name string) (context.Context, Span) {
	return ctx, Span{
		Name:      name,
		StartTime: time.Now(),
		Tags:      make(map[string]string),
		Events:    make([]Event, 0),
	}
}

// EndSpan implements the Tracer interface
func (t *NoopTracer) EndSpan(span Span) {
	span.EndTime = time.Now()
}

// AddEvent implements the Tracer interface
func (t *NoopTracer) AddEvent(span *Span, name string, tags map[string]string) {
	span.Events = append(span.Events, Event{
		Name:      name,
		Timestamp: time.Now(),
		Tags:      tags,
	})
}

// AddTag implements the Tracer interface
func (t *NoopTracer) AddTag(span *Span, key, value string) {
	if span.Tags == nil {
		span.Tags = make(map[string]string)
	}
	span.Tags[key] = value
}

// NewNoopTracer creates a new no-op tracer
func NewNoopTracer() *NoopTracer {
	return &NoopTracer{}
} 