package request

import (
	"context"
	"fmt"
)

type requestInfoCtx struct{}
type requestKeyCtx struct{}
type requestInputCtx struct{}

func BaseContext[T any](ctx context.Context, req *requestInfo, requestKey string, input T) context.Context {
	ctx = context.WithValue(ctx, requestInfoCtx{}, req)
	ctx = context.WithValue(ctx, requestKeyCtx{}, requestKey)
	ctx = context.WithValue(ctx, requestInputCtx{}, input)

	return ctx
}

func GetRequiredScope(ctx context.Context) string {
	if req, ok := ctx.Value(requestInfoCtx{}).(*requestInfo); ok {
		return req.RequiredScope
	}
	return ""
}

func GetRequestInput[T any](ctx context.Context) T {
	if input, ok := ctx.Value(requestInputCtx{}).(T); ok {
		return input
	}
	return *new(T)
}

func GetRequestKey(ctx context.Context) string {
	if key, ok := ctx.Value(requestKeyCtx{}).(string); ok {
		return key
	}
	return ""
}

func GetRoute(ctx context.Context) (string, bool) {
	if req, ok := ctx.Value(requestInfoCtx{}).(*requestInfo); ok {
		return fmt.Sprintf("%s %s", req.Method, req.Path), true
	}
	return "", false
}

// WithRoute returns a context carrying route metadata for middleware that uses GetRoute.
func WithRoute(ctx context.Context, method, path string) context.Context {
	return context.WithValue(ctx, requestInfoCtx{}, &requestInfo{Method: method, Path: path})
}
