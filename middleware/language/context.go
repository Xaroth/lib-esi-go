package language

import (
	"context"

	"github.com/xaroth/lib-esi-go/request"
)

type requestLanguageCtx struct{}

func WithLanguage(language string) request.RequestOption {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, requestLanguageCtx{}, language)
	}
}

func getLanguage(ctx context.Context) (string, bool) {
	language, ok := ctx.Value(requestLanguageCtx{}).(string)
	return language, ok
}
