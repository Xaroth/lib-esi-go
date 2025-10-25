package useragent_test

import (
	"errors"
	"net/http"
	"testing"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/useragent"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		applicationName    string
		applicationVersion string
		applicationContact []string
		expectedUserAgent  string
		expectedError      error
	}{
		{
			name:               "success: application name, version, and contact are present",
			applicationName:    "test-application",
			applicationVersion: "1.0.0",
			applicationContact: []string{"test-contact"},
			expectedUserAgent:  "test-application/1.0.0 (test-contact); " + defaults.UserAgent,
		},
		{
			name:               "success: application name and version are present, contact is not present",
			applicationName:    "test-application",
			applicationVersion: "1.0.0",
			applicationContact: []string{},
			expectedUserAgent:  "test-application/1.0.0; " + defaults.UserAgent,
		},
		{
			name:               "success: application name is present, version and contact are not present",
			applicationName:    "test-application",
			applicationVersion: "",
			applicationContact: []string{},
			expectedUserAgent:  "test-application; " + defaults.UserAgent,
		},
		{
			name:               "success: multiple contacts are present",
			applicationName:    "test-application",
			applicationVersion: "1.0.0",
			applicationContact: []string{"test-contact1", "test-contact2"},
			expectedUserAgent:  "test-application/1.0.0 (test-contact1; test-contact2); " + defaults.UserAgent,
		},
		{
			name:               "success: application name, version, and contact are not present",
			applicationName:    "",
			applicationVersion: "",
			applicationContact: []string{},
			expectedUserAgent:  "Unconfigured; " + defaults.UserAgent,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			expectation := func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("User-Agent") != testCase.expectedUserAgent {
					tb.Fatalf("expected User-Agent header to be '%s', got '%s'", testCase.expectedUserAgent, req.Header.Get("User-Agent"))
				}
			}

			resp, err := useragent.Middleware(testCase.applicationName, testCase.applicationVersion, testCase.applicationContact...)(middleware.NewFakeMiddleware(t, expectation)).RoundTrip(req)

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
