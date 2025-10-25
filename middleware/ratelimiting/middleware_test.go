package ratelimiting_test

import (
	"net/http"
	"testing"

	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/ratelimiting"
	"github.com/xaroth/lib-esi-go/request"
)

type fakeRateLimiter struct {
	scheduled bool
}

func (f *fakeRateLimiter) Schedule(req *http.Request) (func(*http.Response), error) {
	f.scheduled = true
	return func(*http.Response) {}, nil
}

func TestMiddleware_skipsWithoutRoute(t *testing.T) {
	t.Parallel()

	limiter := &fakeRateLimiter{}
	rt := ratelimiting.Middleware(limiter)(middleware.NewFakeMiddleware(t, func(tb testing.TB, req *http.Request) {
		tb.Helper()
	}))

	req, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := rt.RoundTrip(req); err != nil {
		t.Fatal(err)
	}
	if limiter.scheduled {
		t.Fatal("expected rate limiter not to run without route in context")
	}
}

func TestMiddleware_schedulesWithRoute(t *testing.T) {
	t.Parallel()

	limiter := &fakeRateLimiter{}
	rt := ratelimiting.Middleware(limiter)(middleware.NewFakeMiddleware(t, func(tb testing.TB, req *http.Request) {
		tb.Helper()
	}))

	req, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(request.WithRoute(req.Context(), http.MethodGet, "/foo"))

	if _, err := rt.RoundTrip(req); err != nil {
		t.Fatal(err)
	}
	if !limiter.scheduled {
		t.Fatal("expected rate limiter to run when route is in context")
	}
}

func TestMiddleware_nilLimiterPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	_ = ratelimiting.Middleware(nil)
}
