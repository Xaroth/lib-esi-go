package timeout

import (
	"context"
	"net/http"
	"time"

	"github.com/xaroth/lib-esi-go/middleware"
)

// Middleware sets a timeout for the request when configured.
// This middleware is always added to the transport chain.
func Middleware(defaultTimeout time.Duration) middleware.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
			timeout := defaultTimeout
			if value, ok := getTimeout(req.Context()); ok {
				timeout = value
			}

			if timeout <= 0 {
				return next.RoundTrip(req)
			}

			ctx := req.Context()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			req = req.WithContext(ctx)
			return next.RoundTrip(req)
		})
	}
}
