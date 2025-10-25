package ratelimiting

import (
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/authentication"
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
			if info, ok := request.GetRequestInfo(ctx); ok {
				owner := int64(-1)
				if token, ok := authentication.GetToken(ctx); ok {
					owner = token.Owner()
				}

				done, err := rateLimiter.Schedule(ctx, info, owner)
				if err != nil {
					return nil, err
				}

				resp, err := next.RoundTrip(req)
				done(resp)

				return resp, err
			}
			return next.RoundTrip(req)
		})
	}
}
