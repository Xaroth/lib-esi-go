package compatibilitydate_test

import (
	"errors"
	"net/http"
	"testing"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/compatibilitydate"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		compatibilityDate string
		expectation       func(tb testing.TB, req *http.Request)
		expectedError     error
	}{
		{
			name:              "success: compatibility date is present",
			compatibilityDate: "2025-01-01",
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("X-Compatibility-Date") != "2025-01-01" {
					tb.Fatalf("expected X-Compatibility-Date header to be '2025-01-01', got '%s'", req.Header.Get("X-Compatibility-Date"))
				}
			},
		},
		{
			name:              "success: compatibility date falls back to default",
			compatibilityDate: "",
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("X-Compatibility-Date") != defaults.CompatibilityDate {
					tb.Fatalf("expected X-Compatibility-Date header to be empty, got '%s'", req.Header.Get("X-Compatibility-Date"))
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := compatibilitydate.Override(t.Context(), testCase.compatibilityDate)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := compatibilitydate.Middleware(middleware.NewFakeMiddleware(t, testCase.expectation)).RoundTrip(req)

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
