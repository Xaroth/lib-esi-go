package language

import (
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware"
)

// Middleware adds an Accept-Language header to each request when configured.
// This middleware is always added to the transport chain.
func Middleware(defaultLanguage string) middleware.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
			language := defaultLanguage
			if value, ok := getLanguage(req.Context()); ok {
				language = value
			}

			if language == "" {
				return next.RoundTrip(req)
			}

			req = req.Clone(req.Context())
			req.Header.Set("Accept-Language", language)
			return next.RoundTrip(req)
		})
	}
}
