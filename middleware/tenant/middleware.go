package tenant

import (
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware"
)

// Middleware adds an X-Tenant header to each request when configured.
// This middleware is always added to the transport chain.
func Middleware(defaultTenant string) middleware.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
			tenant := defaultTenant
			if value, ok := getTenant(req.Context()); ok {
				tenant = value
			}

			if tenant == "" {
				return next.RoundTrip(req)
			}

			req = req.Clone(req.Context())
			req.Header.Set("X-Tenant", tenant)
			return next.RoundTrip(req)
		})
	}
}
