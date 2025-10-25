package compatibilitydate

import (
	"context"

	"github.com/xaroth/lib-esi-go/request"
)

type requestCompatibilityDateCtx struct{}

func WithCompatibilityDate(compatibilityDate string) request.RequestOption {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, requestCompatibilityDateCtx{}, compatibilityDate)
	}
}

func getCompatibilityDate(ctx context.Context) (string, bool) {
	compatibilityDate, ok := ctx.Value(requestCompatibilityDateCtx{}).(string)
	return compatibilityDate, ok
}
