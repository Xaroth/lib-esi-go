package authentication_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/authentication"
	"github.com/xaroth/lib-esi-go/middleware/authentication/mock"
	"go.uber.org/mock/gomock"
)

//go:generate go run -mod=mod go.uber.org/mock/mockgen -build_flags=--mod=mod -destination=mock/mock_token.go -package=mock github.com/xaroth/lib-esi-go/middleware/authentication Token,RefreshableToken

func TestMiddleware(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	token := mock.NewMockToken(ctrl)
	token.EXPECT().Token().Return("test-token").AnyTimes()

	refreshable := mock.NewMockRefreshableToken(ctrl)
	refreshable.EXPECT().RefreshIfNeeded(gomock.Any()).Return(nil).AnyTimes()
	refreshable.EXPECT().Token().Return("refreshed-token").AnyTimes()

	refreshError := errors.New("refresh error")
	errorRefresh := mock.NewMockRefreshableToken(ctrl)
	errorRefresh.EXPECT().RefreshIfNeeded(gomock.Any()).Return(refreshError).AnyTimes()
	errorRefresh.EXPECT().Token().Times(0)

	testCases := []struct {
		name          string
		token         authentication.Token
		expectation   func(tb testing.TB, req *http.Request)
		expectedError error
	}{
		{
			name:  "success: token is present",
			token: token,
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("Authorization") != "Bearer test-token" {
					tb.Fatalf("expected authorization header to be 'Bearer test-token', got '%s'", req.Header.Get("Authorization"))
				}
			},
		},
		{
			name:  "success: token is not present",
			token: nil,
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("Authorization") != "" {
					tb.Fatalf("expected authorization header to be empty, got '%s'", req.Header.Get("Authorization"))
				}
			},
		},
		{
			name:  "success: token is refreshable",
			token: refreshable,
			expectation: func(tb testing.TB, req *http.Request) {
				tb.Helper()

				if req.Header.Get("Authorization") != "Bearer refreshed-token" {
					tb.Fatalf("expected authorization header to be 'Bearer refreshed-token', got '%s'", req.Header.Get("Authorization"))
				}
			},
		},
		{
			name:          "error: refresh error",
			token:         errorRefresh,
			expectedError: refreshError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := authentication.WithToken(t.Context(), testCase.token)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := authentication.Middleware(middleware.NewFakeMiddleware(t, testCase.expectation)).RoundTrip(req)

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
