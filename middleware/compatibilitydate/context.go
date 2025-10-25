package compatibilitydate

import (
	"context"

	defaults "github.com/xaroth/lib-esi-go"
)

type requestCompatibilityDateCtx struct{}

func Override(ctx context.Context, compatibilityDate string) context.Context {
	return context.WithValue(ctx, requestCompatibilityDateCtx{}, compatibilityDate)
}

// Returns the compatibility date to use for the request.
// If no override is set, the default compatibility date is returned.
func GetCompatibilityDate(ctx context.Context) string {
	if compatibilityDate, ok := ctx.Value(requestCompatibilityDateCtx{}).(string); ok && compatibilityDate != "" {
		return compatibilityDate
	}

	return defaults.CompatibilityDate
}
