package tenant

import (
	"context"

	"github.com/xaroth/lib-esi-go/request"
)

type requestTenantCtx struct{}

func WithTenant(tenant string) request.RequestOption {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, requestTenantCtx{}, tenant)
	}
}

func getTenant(ctx context.Context) (string, bool) {
	tenant, ok := ctx.Value(requestTenantCtx{}).(string)
	return tenant, ok
}
