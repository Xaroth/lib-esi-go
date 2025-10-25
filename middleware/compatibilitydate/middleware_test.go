package compatibilitydate_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/compatibilitydate"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	const defaultCompatibilityDate = "2026-02-01"

	testCases := []struct {
		name              string
		useContext        bool
		compatibilityDate string
		expectation       func(tb testing.TB, req *http.Request)
		expectedError     error
	}{
		{
			name:              "success: compatibility date is overridden in context",
			useContext:        true,
			compatibilityDate: "2025-01-01",
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("X-Compatibility-Date") != "2025-01-01" {
					tb.Fatalf("expected X-Compatibility-Date header to be '2025-01-01', got '%s'", req.Header.Get("X-Compatibility-Date"))
				}
			},
		},
		{
			name:       "success: compatibility date is not present in context, default compatibility date is used",
			useContext: false,
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("X-Compatibility-Date") != defaultCompatibilityDate {
					tb.Fatalf("expected X-Compatibility-Date header to be '%s', got '%s'", defaultCompatibilityDate, req.Header.Get("X-Compatibility-Date"))
				}
			},
		},
		{
			name:              "success: compatibility date is disabled for the request",
			useContext:        true,
			compatibilityDate: "",
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("X-Compatibility-Date") != "" {
					tb.Fatalf("expected X-Compatibility-Date header to be absent, got '%s'", req.Header.Get("X-Compatibility-Date"))
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if testCase.useContext {
				ctx = compatibilitydate.WithCompatibilityDate(testCase.compatibilityDate)(ctx)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := compatibilitydate.Middleware(defaultCompatibilityDate)(middleware.NewFakeMiddleware(t, testCase.expectation)).RoundTrip(req)

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

func TestMiddleware_InvalidCompatibilityDate(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got nil")
		}
	}()
	compatibilitydate.Middleware("invalid")
	t.Fatalf("expected panic, got nil")
}
