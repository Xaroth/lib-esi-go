package compatibilitydate

import (
	"fmt"
	"net/http"
	"time"

	"github.com/xaroth/lib-esi-go/middleware"
)

// Middleware adds a X-Compatibility-Date header to each request when configured.
// This middleware is always added to the transport chain.
func Middleware(defaultCompatibilityDate string) middleware.Middleware {
	if defaultCompatibilityDate != "" {
		if _, err := time.Parse(time.DateOnly, defaultCompatibilityDate); err != nil {
			panic(fmt.Errorf("invalid compatibility date: %w", err))
		}
	}

	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
			compatibilityDate := defaultCompatibilityDate
			if value, ok := getCompatibilityDate(req.Context()); ok {
				compatibilityDate = value
			}

			if compatibilityDate == "" {
				return next.RoundTrip(req)
			}

			req = req.Clone(req.Context())
			req.Header.Set("X-Compatibility-Date", compatibilityDate)
			return next.RoundTrip(req)
		})
	}
}
