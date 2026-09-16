package request_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/xaroth/lib-esi-go/request"
)

type ctxKey struct{}

type senderFunc func(req *http.Request) (*http.Response, error)

func (f senderFunc) Do(req *http.Request) (*http.Response, error) { return f(req) }

// Request options must end up on the outgoing request's context; that is how
// middleware such as authentication finds the token.
func TestRequestOptionsReachContext(t *testing.T) {
	t.Parallel()

	var seen any
	sender := senderFunc(func(req *http.Request) (*http.Response, error) {
		seen = req.Context().Value(ctxKey{})
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader(`{}`)),
		}, nil
	})

	get := request.CreateStatic[map[string]any](http.MethodGet, "/status")
	withValue := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxKey{}, "present")
	}

	if _, err := get(context.Background(), sender, withValue); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if seen != "present" {
		t.Fatalf("request option was not applied to the context, got %v", seen)
	}
}
