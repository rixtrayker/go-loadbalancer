package tracer

import (
	"context"
	"time"
)

// Span represents a tracing span
type Span struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Tags      map[string]string
	Events    []Event
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