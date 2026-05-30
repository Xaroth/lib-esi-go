package tier_test

import (
	"errors"
	"net/http"
	"testing"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/tier"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	const defaultTier = defaults.Tier

	testCases := []struct {
		name          string
		tier          string
		initialURL    string
		expectation   func(tb testing.TB, req *http.Request)
		expectedError error
	}{
		{
			name:       "success: default tier does not modify URL host",
			tier:       defaultTier,
			initialURL: "https://original.example.com/path",
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.URL.Host != "original.example.com" {
					tb.Fatalf("expected URL host to be 'original.example.com', got '%s'", req.URL.Host)
				}
			},
		},
		{
			name:       "success: test tier overrides URL host",
			tier:       "test",
			initialURL: "https://original.example.com/path",
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				expected := defaults.TieredHosts["test"].Host
				if req.URL.Host != expected {
					tb.Fatalf("expected URL host to be '%s', got '%s'", expected, req.URL.Host)
				}
			},
		},
		{
			name:       "success: dev tier overrides URL host",
			tier:       "dev",
			initialURL: "https://original.example.com/path",
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				expected := defaults.TieredHosts["dev"].Host
				if req.URL.Host != expected {
					tb.Fatalf("expected URL host to be '%s', got '%s'", expected, req.URL.Host)
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testCase.initialURL, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := tier.Middleware(testCase.tier)(middleware.NewFakeMiddleware(t, testCase.expectation)).RoundTrip(req)

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
