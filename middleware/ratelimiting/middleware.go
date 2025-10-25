package ratelimiting

import (
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/request"
)

// Middleware automatically delays requests to ensure rate limits are respected.
// This middleware is opt-in, and is not enabled by default.
func Middleware(rateLimiter RateLimiter) middleware.Middleware {
	if rateLimiter == nil {
		panic("no rate limiter backend provided")
	}

	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
			ctx := req.Context()

			if _, ok := request.GetRoute(ctx); !ok {
				// If no route is set, we are not processing an ESI request, skip rate limiting.
				return next.RoundTrip(req)
			}

			done, err := rateLimiter.Schedule(req)
			if err != nil {
				return nil, err
			}

			resp, err := next.RoundTrip(req)
			done(resp)

			return resp, err
		})
	}
}
