package tenant_test

import (
	"errors"
	"net/http"
	"testing"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/tenant"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	const defaultTenant = defaults.Tenant

	testCases := []struct {
		name          string
		useContext    bool
		tenant        string
		expectation   func(tb testing.TB, req *http.Request)
		expectedError error
	}{
		{
			name:       "success: tenant is overridden in context",
			useContext: true,
			tenant:     "singularity",
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("X-Tenant") != "singularity" {
					tb.Fatalf("expected X-Tenant header to be 'singularity', got '%s'", req.Header.Get("X-Tenant"))
				}
			},
		},
		{
			name:       "success: tenant is not present in context, default tenant is used",
			useContext: false,
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("X-Tenant") != defaultTenant {
					tb.Fatalf("expected X-Tenant header to be '%s', got '%s'", defaultTenant, req.Header.Get("X-Tenant"))
				}
			},
		},
		{
			name:       "success: tenant is disabled for the request",
			useContext: true,
			tenant:     "",
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("X-Tenant") != "" {
					tb.Fatalf("expected X-Tenant header to be absent, got '%s'", req.Header.Get("X-Tenant"))
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if testCase.useContext {
				ctx = tenant.WithTenant(testCase.tenant)(ctx)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := tenant.Middleware(defaultTenant)(middleware.NewFakeMiddleware(t, testCase.expectation)).RoundTrip(req)

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
