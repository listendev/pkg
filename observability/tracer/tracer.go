package tracer

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ExporterBuilder(ctx context.Context, addr string, filePath string) (sdktrace.SpanExporter, error) {
	if addr == "" {
		if filePath == "" {
			return nil, fmt.Errorf("no address or file path provided, cannot make a span exporter")
		}
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}

		return NewStdoutExporter(f)
	}

	return NewGRPCExporter(ctx, addr)
}

func NewStdoutExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithWriter(w))
}

func NewGRPCExporter(ctx context.Context, addr string) (sdktrace.SpanExporter, error) {
	// todo: figure out if we need tls here as we probably will use otel-collector
	// which runs on the same host anyways.
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	return otlptracegrpc.New(timeoutCtx, otlptracegrpc.WithGRPCConn(conn))
}

func NewTraceProvider(ctx context.Context, exp sdktrace.SpanExporter, serviceName string) (*sdktrace.TracerProvider, error) {
	var err error

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}

type ctxKey struct{}

func FromContext(ctx context.Context) trace.Tracer {
	t, _ := ctx.Value(ctxKey{}).(trace.Tracer)

	return t
}

func NewContext(parent context.Context, t trace.Tracer) context.Context {
	return context.WithValue(parent, ctxKey{}, t)
}
