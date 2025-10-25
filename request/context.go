package request

import (
	"context"
	"net/url"
	"time"

	defaults "github.com/xaroth/lib-esi-go"
)

type requestStartCtx struct{}
type requestInfoCtx struct{}
type requestInputCtx struct{}
type requestHostCtx struct{}
type requestKeyCtx struct{}

func OverrideHost(ctx context.Context, host *url.URL) context.Context {
	return context.WithValue(ctx, requestHostCtx{}, host)
}

// GetHost extracts the host from the context.
func GetHost(ctx context.Context) *url.URL {
	if host, ok := ctx.Value(requestHostCtx{}).(*url.URL); ok {
		return host
	}

	return defaults.Host
}

func WithRequestContext[T any](
	ctx context.Context,
	info RequestInfo,
	key string,
	input *T,
) context.Context {
	ctx = context.WithValue(ctx, requestStartCtx{}, time.Now())
	ctx = context.WithValue(ctx, requestInfoCtx{}, info)
	ctx = context.WithValue(ctx, requestKeyCtx{}, key)
	ctx = context.WithValue(ctx, requestInputCtx{}, input)
	return ctx
}

// GetRequestStart extracts the request start time from the context.
func GetRequestStart(ctx context.Context) (time.Time, bool) {
	start, ok := ctx.Value(requestStartCtx{}).(time.Time)
	return start, ok
}

// GetRequestInfo returns the RequestInfo from the context.
func GetRequestInfo(ctx context.Context) (RequestInfo, bool) {
	info, ok := ctx.Value(requestInfoCtx{}).(RequestInfo)
	return info, ok
}

// GetRequestInput returns the request input from the context.
func GetRequestInput[T any](ctx context.Context) (T, bool) {
	val, ok := ctx.Value(requestInputCtx{}).(T)
	return val, ok
}

// GetRequestKey returns the request key from the context.
// This key is a combination of all path, query and header parameters, sorted.
func GetRequestKey(ctx context.Context) string {
	if key, ok := ctx.Value(requestKeyCtx{}).(string); ok {
		return key
	}
	return ""
}
