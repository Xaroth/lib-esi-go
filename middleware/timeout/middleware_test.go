package timeout_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/timeout"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	const defaultTimeout = defaults.RequestTimeout

	testCases := []struct {
		name          string
		useContext    bool
		timeout       time.Duration
		expectation   func(tb testing.TB, req *http.Request)
		expectedError error
	}{
		{
			name:       "success: timeout is overridden in context",
			useContext: true,
			timeout:    5 * time.Second,
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if deadline, ok := req.Context().Deadline(); !ok || deadline.IsZero() {
					tb.Fatalf("expected deadline to be set")
				}
			},
		},
		{
			name:       "success: timeout is not present in context, default timeout is used",
			useContext: false,
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if deadline, ok := req.Context().Deadline(); !ok || deadline.IsZero() {
					tb.Fatalf("expected deadline to be set")
				}
			},
		},
		{
			name:       "success: timeout is disabled for the request",
			useContext: true,
			timeout:    0,
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if _, ok := req.Context().Deadline(); ok {
					tb.Fatalf("expected deadline to be absent")
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if testCase.useContext {
				ctx = timeout.WithTimeout(testCase.timeout)(ctx)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := timeout.Middleware(defaultTimeout)(middleware.NewFakeMiddleware(t, testCase.expectation)).RoundTrip(req)

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("failed to round trip request: %v", err)
				}
				if !errors.Is(err, testCase.expectedError) {
					t.Fatalf("expected error %v, got %v", testCase.expectedError, err)
				}
				return
			}

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}
		})
	}
}
