package timeout

import (
	"context"
	"time"

	defaults "github.com/xaroth/lib-esi-go"
)

type requestTimeoutCtx struct{}

func Override(ctx context.Context, timeout time.Duration) context.Context {
	return context.WithValue(ctx, requestTimeoutCtx{}, timeout)
}

// Returns the timeout to use for the request.
// If no override is set, the default timeout is returned.
func GetTimeout(ctx context.Context) time.Duration {
	if timeout, ok := ctx.Value(requestTimeoutCtx{}).(time.Duration); ok && timeout.Seconds() > 0.0 {
		return timeout
	}

	return defaults.RequestTimeout
}
