package tracer

import (
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func NewNoopTracerProvider() trace.TracerProvider {
	return noop.NewTracerProvider()
}
