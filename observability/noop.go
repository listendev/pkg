package observability

import (
	"context"

	"github.com/listendev/pkg/observability/logger"
	"github.com/listendev/pkg/observability/threadid"
	"github.com/listendev/pkg/observability/tracer"
	"go.uber.org/zap"
)

func NewNopContext() context.Context {
	traceCtx := tracer.NewContext(context.Background(), tracer.NoOpTracer{})
	logCtx := logger.NewContext(traceCtx, zap.NewNop())
	threadIDCtx := threadid.NewContext(logCtx, threadid.ThreadID(0))

	return threadIDCtx
}
