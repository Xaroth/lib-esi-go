package authentication_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/authentication"
	"github.com/xaroth/lib-esi-go/middleware/authentication/mock"
	"github.com/xaroth/lib-esi-go/request"
	"go.uber.org/mock/gomock"
)

//go:generate go run -mod=mod go.uber.org/mock/mockgen -build_flags=--mod=mod -destination=mock/mock_token.go -package=mock github.com/xaroth/lib-esi-go/middleware/authentication Token,RefreshableToken,ScopedToken

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

	scoped := mock.NewMockScopedToken(ctrl)
	scoped.EXPECT().Scopes().Return([]string{"scope1", "scope2"}).AnyTimes()
	scoped.EXPECT().Token().Return("test-token").AnyTimes()

	testCases := []struct {
		name           string
		token          authentication.Token
		requiredScopes []string
		expectation    func(tb testing.TB, req *http.Request)
		expectedError  error
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
		{
			name:           "error: missing token but scopes are required",
			token:          nil,
			requiredScopes: []string{"missing"},
			expectedError:  authentication.ErrMissingToken,
		},
		{
			name:           "success: token is scoped and scopes are present",
			token:          scoped,
			requiredScopes: []string{"scope1"},
		},
		{
			name:           "success: token is scoped and at least one scope is present",
			token:          scoped,
			requiredScopes: []string{"scope2", "scope3"},
		},
		{
			name:           "error: missing scopes but token is present",
			token:          scoped,
			requiredScopes: []string{"missing"},
			expectedError:  authentication.ErrMissingScopes,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			reqInfo := request.FakeRequestInfo(
				http.MethodGet,
				"/",
				request.WithRequiredScope(testCase.requiredScopes...),
			)

			ctx := request.BaseContext[any](t.Context(), reqInfo, "test", nil)

			ctx = authentication.WithToken(testCase.token)(ctx)
			req, err := http.NewRequestWithContext(ctx, reqInfo.Method, reqInfo.Path, nil)
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
