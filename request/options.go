package request

import "context"

type CreateOption func(*requestInfo)

type RequestOption func(context.Context) context.Context

func WithRequiredScope(scope ...string) CreateOption {
	return func(info *requestInfo) {
		info.RequiredScope = append(info.RequiredScope, scope...)
	}
}
