package authentication

import "context"

type requestTokenCtx struct{}

func WithToken(ctx context.Context, token Token) context.Context {
	return context.WithValue(ctx, requestTokenCtx{}, token)
}

// GetToken extracts the Token from the context.
// Returns the Token and true if found, or nil and false if not found.
func GetToken(ctx context.Context) (Token, bool) {
	token, ok := ctx.Value(requestTokenCtx{}).(Token)
	return token, ok
}
