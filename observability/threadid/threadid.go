package threadid

import "context"

type ThreadID int
type ctxKey struct{}

func FromContext(ctx context.Context) ThreadID {
	id, ok := ctx.Value(ctxKey{}).(ThreadID)
	if !ok {
		return 0
	}

	return id
}

func NewContext(parent context.Context, id ThreadID) context.Context {
	return context.WithValue(parent, ctxKey{}, id)
}
