package logger

import (
	"context"
	"strconv"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

func NewContext(parent context.Context, z *zap.Logger) context.Context {
	return context.WithValue(parent, ctxKey{}, z)
}

func New(level zapcore.Level, json bool, options ...zap.Option) (*zap.Logger, error) {
	var cfg zap.Config

	cfg = zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if json {
		cfg = zap.Config{
			Level:       zap.NewAtomicLevelAt(level),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}
	return cfg.Build(options...)
}

func convertIDToDatadogFormat(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)

}
func FromContext(ctx context.Context) *zap.Logger {
	childLogger, _ := ctx.Value(ctxKey{}).(*zap.Logger)

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return childLogger
	}
	if traceID := span.SpanContext().TraceID(); traceID.IsValid() {
		childLogger = childLogger.With(
			zap.String("trace_id", traceID.String()),
			zap.String("dd.trace_id", convertIDToDatadogFormat(traceID.String())),
		)
	}

	if spanID := span.SpanContext().SpanID(); spanID.IsValid() {
		childLogger = childLogger.With(
			zap.String("span_id", spanID.String()),
			zap.String("dd.span_id", convertIDToDatadogFormat(spanID.String())),
		)
	}
	return childLogger
}
