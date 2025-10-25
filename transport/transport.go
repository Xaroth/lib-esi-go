package transport

import (
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/authentication"
	"github.com/xaroth/lib-esi-go/middleware/compatibilitydate"
	"github.com/xaroth/lib-esi-go/middleware/timeout"
	"github.com/xaroth/lib-esi-go/middleware/useragent"
)

type transportChain struct {
	base        http.RoundTripper
	middlewares []middleware.Middleware
}

func (c *transportChain) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := c.base

	// Reverse the middleware chain to ensure that the middleware is applied in the correct order.
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		rt = c.middlewares[i](rt)
	}

	return rt.RoundTrip(req)

}

func newTransportChain(base http.RoundTripper, middlewares ...middleware.Middleware) *transportChain {
	return &transportChain{
		base:        base,
		middlewares: middlewares,
	}
}

// NewESITransport constructs a transport chain with the default middlewares.
// This is the recommended way to construct a transport chain for ESI requests.
func NewESITransport(applicationName, applicationVersion string, contact []string, opts ...Option) http.RoundTripper {
	chain := newTransportChain(
		http.DefaultTransport,
		timeout.Middleware,
		useragent.Middleware(applicationName, applicationVersion, contact...),
		compatibilitydate.Middleware,
		authentication.Middleware,
	)

	for _, opt := range opts {
		opt(chain)
	}

	return chain
}

// NewTransportChain constructs an empty transport chain with only the given middlewares.
// This can be used in advanced cases where the default middlewares are not needed.
func NewTransportChain(base http.RoundTripper, middlewares ...middleware.Middleware) http.RoundTripper {
	return newTransportChain(base, middlewares...)
}
