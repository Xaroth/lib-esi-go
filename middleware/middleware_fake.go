package middleware

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

// NewFakeMiddleware creates a downstream http.RoundTripper that always returns a 200 OK response,
// to aid with testing middleware.
func NewFakeMiddleware(tb testing.TB, validator func(tb testing.TB, req *http.Request)) MiddlewareFunc {
	tb.Helper()

	return MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
		tb.Helper()

		validator(tb, req)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})
}
