package language_test

import (
	"errors"
	"net/http"
	"testing"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/language"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	const defaultLanguage = defaults.Language

	testCases := []struct {
		name          string
		useContext    bool
		language      string
		expectation   func(tb testing.TB, req *http.Request)
		expectedError error
	}{
		{
			name:       "success: language is overridden in context",
			useContext: true,
			language:   "de",
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("Accept-Language") != "de" {
					tb.Fatalf("expected Accept-Language header to be 'de', got '%s'", req.Header.Get("Accept-Language"))
				}
			},
		},
		{
			name:       "success: language is not present in context, default language is used",
			useContext: false,
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("Accept-Language") != defaultLanguage {
					tb.Fatalf("expected Accept-Language header to be '%s', got '%s'", defaultLanguage, req.Header.Get("Accept-Language"))
				}
			},
		},
		{
			name:       "success: language is disabled for the request",
			useContext: true,
			language:   "",
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("Accept-Language") != "" {
					tb.Fatalf("expected Accept-Language header to be absent, got '%s'", req.Header.Get("Accept-Language"))
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if testCase.useContext {
				ctx = language.WithLanguage(testCase.language)(ctx)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := language.Middleware(defaultLanguage)(middleware.NewFakeMiddleware(t, testCase.expectation)).RoundTrip(req)

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
