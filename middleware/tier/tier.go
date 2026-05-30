package tier

import (
	"net/http"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/middleware"
)

func Middleware(tier string) middleware.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
			if tier != defaults.Tier {
				req = req.Clone(req.Context())

				host := defaults.TieredHosts[tier]
				req.URL.Host = host.Host
			}

			return next.RoundTrip(req)
		})
	}
}
