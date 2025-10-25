package authentication

import "context"

type Token interface {
	// The identifier of the entity that this token belongs to.
	Owner() int64

	// The current token
	Token() string
}

type RefreshableToken interface {
	Token

	// Called to check if the token needs to be refreshed, and if so, refresh it.
	RefreshIfNeeded(ctx context.Context) error
}
