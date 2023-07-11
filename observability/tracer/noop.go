package tracer

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type NoOpTracerProvider struct{}

var _ trace.TracerProvider = NoOpTracerProvider{}

// Tracer returns noop implementation of Tracer.
func (p NoOpTracerProvider) Tracer(string, ...trace.TracerOption) trace.Tracer {
	return NoOpTracer{}
}

type NoOpTracer struct{}

var _ trace.Tracer = NoOpTracer{}

// Start carries forward a non-recording Span, if one is present in the context, otherwise it
// creates a no-op Span.
func (t NoOpTracer) Start(ctx context.Context, _ string, _ ...trace.SpanStartOption) (context.Context, trace.Span) {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		span = NoOpSpan{}
	}

	return trace.ContextWithSpan(ctx, span), span
}

// NoOpSpan is an implementation of Span that preforms no operations.
type NoOpSpan struct{}

var _ trace.Span = NoOpSpan{}

// SpanContext returns an empty span context.
func (NoOpSpan) SpanContext() trace.SpanContext { return trace.SpanContext{} }

// IsRecording always returns false.
func (NoOpSpan) IsRecording() bool { return false }

// SetStatus does nothing.
func (NoOpSpan) SetStatus(codes.Code, string) {}

// SetError does nothing.
func (NoOpSpan) SetError(bool) {}

// SetAttributes does nothing.
func (NoOpSpan) SetAttributes(...attribute.KeyValue) {}

// End does nothing.
func (NoOpSpan) End(...trace.SpanEndOption) {}

// RecordError does nothing.
func (NoOpSpan) RecordError(error, ...trace.EventOption) {}

// AddEvent does nothing.
func (NoOpSpan) AddEvent(string, ...trace.EventOption) {}

// SetName does nothing.
func (NoOpSpan) SetName(string) {}

// TracerProvider returns a no-op TracerProvider.
func (NoOpSpan) TracerProvider() trace.TracerProvider { return NoOpTracerProvider{} }
