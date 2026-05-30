package transport

import (
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware"
)

type Option func(*transportChain)

func WithTransport(transport http.RoundTripper) Option {
	return func(c *transportChain) {
		c.base = transport
	}
}

func WithMiddleware(middleware middleware.Middleware) Option {
	return func(c *transportChain) {
		c.middlewares = append(c.middlewares, middleware)
	}
}

func WithMiddlewares(middlewares ...middleware.Middleware) Option {
	return func(c *transportChain) {
		c.middlewares = append(c.middlewares, middlewares...)
	}
}

func WithTier(tier string) Option {
	return func(c *transportChain) {
		c.defaultTier = tier
	}
}

func WithTenant(tenant string) Option {
	return func(c *transportChain) {
		c.defaultTenant = tenant
	}
}
