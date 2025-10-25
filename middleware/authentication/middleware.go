package authentication

import (
	"fmt"
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware"
)

// Middleware adds a Authorization header to each request.
// This middleware is always added to the transport chain.
func Middleware(next http.RoundTripper) http.RoundTripper {
	return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
		ctx := req.Context()
		if token, ok := GetToken(ctx); ok {
			if refreshable, ok := token.(RefreshableToken); ok {
				if err := refreshable.RefreshIfNeeded(ctx); err != nil {
					return nil, err
				}
			}

			req = req.Clone(ctx)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token()))
		}
		return next.RoundTrip(req)
	})
}
