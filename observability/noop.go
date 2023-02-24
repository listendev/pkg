package observability

import (
	"context"

	"github.com/garnet-org/pkg/observability/logger"
	"github.com/garnet-org/pkg/observability/threadid"
	"github.com/garnet-org/pkg/observability/tracer"
	"go.uber.org/zap"
)

func NewNopContext() context.Context {
	traceCtx := tracer.NewContext(context.Background(), tracer.NoOpTracer{})
	logCtx := logger.NewContext(traceCtx, zap.NewNop())
	threadIdCtx := threadid.NewContext(logCtx, threadid.ThreadID(0))
	return threadIdCtx
}
