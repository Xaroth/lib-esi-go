package compatibilitydate

import (
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware"
)

// Middleware adds a X-Compatibility-Date header to each request.
// This middleware is always added to the transport chain.
func Middleware(next http.RoundTripper) http.RoundTripper {
	return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
		ctx := req.Context()

		compatibilityDate := GetCompatibilityDate(ctx)

		req = req.Clone(ctx)
		req.Header.Set("X-Compatibility-Date", compatibilityDate)
		return next.RoundTrip(req)
	})
}
