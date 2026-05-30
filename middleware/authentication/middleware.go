package authentication

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/request"
)

var (
	ErrMissingAuthentication = errors.New("missing authentication")
	ErrMissingToken          = errors.New("missing token")
	ErrMissingScopes         = errors.New("missing scopes")
)

// Middleware adds a Authorization header to each request.
// This middleware is always added to the transport chain.
func Middleware(next http.RoundTripper) http.RoundTripper {
	return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
		ctx := req.Context()

		hasToken := false
		var hasScopes []string

		if token, ok := GetToken(ctx); ok {
			hasToken = true

			if refreshable, ok := token.(RefreshableToken); ok {
				if err := refreshable.RefreshIfNeeded(ctx); err != nil {
					return nil, err
				}
			}

			if scoped, ok := token.(ScopedToken); ok {
				hasScopes = scoped.Scopes()
			}

			req = req.Clone(ctx)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token()))
		}

		requiredScopes := request.GetRequiredScope(ctx)
		if len(requiredScopes) > 0 {
			if !hasToken {
				return nil, errors.Join(ErrMissingAuthentication, ErrMissingToken)
			}

			if hasScopes != nil {
				if !hasAnyScope(hasScopes, requiredScopes) {
					return nil, errors.Join(ErrMissingAuthentication, ErrMissingScopes)
				}
			}
		}

		return next.RoundTrip(req)
	})
}

func hasAnyScope(scopes []string, requiredScopes []string) bool {
	for _, scope := range requiredScopes {
		if slices.Contains(scopes, scope) {
			return true
		}
	}
	return false
}
