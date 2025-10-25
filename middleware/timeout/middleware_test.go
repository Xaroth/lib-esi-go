package timeout_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/timeout"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		timeout       time.Duration
		expectation   func(tb testing.TB, req *http.Request)
		expectedError error
	}{
		{
			name:    "success: timeout is present",
			timeout: 10 * time.Second,
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if deadline, ok := req.Context().Deadline(); !ok || deadline.IsZero() {
					tb.Fatalf("expected deadline to be set")
				}
			},
		},
		{
			name:    "success: timeout is not present, default timeout is used",
			timeout: 0 * time.Second,
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if deadline, ok := req.Context().Deadline(); !ok || deadline.IsZero() {
					tb.Fatalf("expected deadline to be set")
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := timeout.Override(t.Context(), testCase.timeout)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := timeout.Middleware(middleware.NewFakeMiddleware(t, testCase.expectation)).RoundTrip(req)

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
