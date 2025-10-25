package cache

import (
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware"

	"github.com/bartventer/httpcache"
)

func Middleware(path string) middleware.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return httpcache.NewTransport(
			CacheDriverName+"://"+path,
			httpcache.WithUpstream(next),
		)
	}
}
