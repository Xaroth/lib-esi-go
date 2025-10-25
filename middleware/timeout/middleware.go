package timeout

import (
	"context"
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware"
)

// Middleware sets a timeout for the request.
// This middleware is always added to the transport chain.
func Middleware(next http.RoundTripper) http.RoundTripper {
	return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
		ctx := req.Context()

		ctx, cancel := context.WithTimeout(ctx, GetTimeout(ctx))
		defer cancel()

		req = req.WithContext(ctx)

		return next.RoundTrip(req)
	})
}
