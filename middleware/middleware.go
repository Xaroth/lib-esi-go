package middleware

import (
	"net/http"
)

type Middleware func(http.RoundTripper) http.RoundTripper

// MiddlewareFunc is similar to http.HandlerFunc, an adapter for taking a normal function and turning it
// into a http.RoundTripper.
type MiddlewareFunc func(r *http.Request) (*http.Response, error)

func (mw MiddlewareFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return mw(req)
}

// Noop is a middleware that does nothing.
func Noop(next http.RoundTripper) http.RoundTripper {
	return next
}
