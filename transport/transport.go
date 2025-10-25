package transport

import (
	"net/http"
	"time"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/authentication"
	"github.com/xaroth/lib-esi-go/middleware/compatibilitydate"
	"github.com/xaroth/lib-esi-go/middleware/language"
	"github.com/xaroth/lib-esi-go/middleware/tenant"
	"github.com/xaroth/lib-esi-go/middleware/timeout"
	"github.com/xaroth/lib-esi-go/middleware/useragent"
)

type transportChain struct {
	base        http.RoundTripper
	middlewares []middleware.Middleware
	chain       http.RoundTripper

	defaultTenant            string
	defaultLanguage          string
	defaultTimeout           time.Duration
	defaultCompatibilityDate string
}

func (c *transportChain) RoundTrip(req *http.Request) (*http.Response, error) {
	return c.chain.RoundTrip(req)
}

// newTransportChain constructs a new transport chain with the given base transport.
// Note that at this point the chain is not assembled, and any public constructors should call
// assemble() to ensure the chain is properly assembled.
func newTransportChain(base http.RoundTripper) *transportChain {
	return &transportChain{
		base:        base,
		middlewares: make([]middleware.Middleware, 0),
		chain:       base,

		defaultTenant:            defaults.Tenant,
		defaultLanguage:          defaults.Language,
		defaultTimeout:           defaults.RequestTimeout,
		defaultCompatibilityDate: defaults.CompatibilityDate,
	}
}

func (c *transportChain) assemble() {
	// Reverse the middleware chain to ensure that the middleware are applied in the correct order.
	// The last middleware in the chain is the first middleware to be called.
	// This ensures that the default transport (which does the actual request) is the last middleware to be called,
	// allowing all other middlewares to modify/inspect/abort the request.
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		c.chain = c.middlewares[i](c.chain)
	}

}

// New constructs a transport chain with the default middlewares.
// This is the recommended way to construct a transport chain for ESI requests.
func New(applicationName, applicationVersion string, contact []string, compatibilityDate string, opts ...Option) http.RoundTripper {
	chain := newTransportChain(http.DefaultTransport)
	defer chain.assemble()

	for _, opt := range opts {
		opt(chain)
	}

	// These middlewares are always added, and options should not be able to override them.
	middlewares := []middleware.Middleware{
		timeout.Middleware(chain.defaultTimeout),
		useragent.Middleware(applicationName, applicationVersion, contact...),
		compatibilitydate.Middleware(compatibilityDate),
		language.Middleware(chain.defaultLanguage),
		tenant.Middleware(chain.defaultTenant),
		authentication.Middleware,
	}
	middlewares = append(middlewares, chain.middlewares...)
	chain.middlewares = middlewares

	return chain
}

// NewChain constructs an empty transport chain with only the given middlewares.
// This can be used in advanced cases where the default middlewares are not needed.
func NewChain(base http.RoundTripper, middlewares ...middleware.Middleware) http.RoundTripper {
	chain := newTransportChain(base)
	chain.middlewares = middlewares
	defer chain.assemble()

	return chain
}
