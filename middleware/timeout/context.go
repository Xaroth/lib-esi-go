package timeout

import (
	"context"
	"time"

	"github.com/xaroth/lib-esi-go/request"
)

type requestTimeoutCtx struct{}

func WithTimeout(timeout time.Duration) request.RequestOption {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, requestTimeoutCtx{}, timeout)
	}
}

func getTimeout(ctx context.Context) (time.Duration, bool) {
	timeout, ok := ctx.Value(requestTimeoutCtx{}).(time.Duration)
	return timeout, ok
}
