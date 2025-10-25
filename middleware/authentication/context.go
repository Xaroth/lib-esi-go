package authentication

import (
	"context"

	"github.com/xaroth/lib-esi-go/request"
)

type requestTokenCtx struct{}

func WithToken(token Token) request.RequestOption {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, requestTokenCtx{}, token)
	}
}

// GetToken extracts the Token from the context.
// Returns the Token and true if found, or nil and false if not found.
func GetToken(ctx context.Context) (Token, bool) {
	token, ok := ctx.Value(requestTokenCtx{}).(Token)
	return token, ok
}
