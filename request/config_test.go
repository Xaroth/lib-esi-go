package request_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/internal/config"
	"github.com/xaroth/lib-esi-go/request"
	"github.com/xaroth/lib-esi-go/request/mock"
	"go.uber.org/mock/gomock"
)

//go:generate go run -mod=mod go.uber.org/mock/mockgen -build_flags=--mod=mod -destination mock/mock_defaults_provider.go -package mock github.com/xaroth/lib-esi-go/request DefaultsProvider

type expectConfig func(config *request.Config)

func newMockToken(t *testing.T, ctrl *gomock.Controller, output string) request.Token {
	t.Helper()

	token := mock.NewMockToken(ctrl)
	token.EXPECT().Token().Return(output).Times(1)
	token.EXPECT().Owner().Return(int64(1)).AnyTimes()
	return token
}

func TestNewConfig(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	otherUrl, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testCases := []struct {
		name         string
		setup        []request.Option
		expectations []expectConfig
	}{
		{
			name: "can set language",
			setup: []request.Option{
				request.WithLanguage("fr"),
			},
			expectations: []expectConfig{
				func(config *request.Config) {
					if config.Language() != "fr" {
						t.Errorf("expected language to be fr, got %s", config.Language())
					}
				},
			},
		},
		{
			name: "can set timeouts",
			setup: []request.Option{
				request.WithScheduleTimeout(-1 * time.Second),
				request.WithRequestTimeout(-1 * time.Second),
			},
			expectations: []expectConfig{
				func(config *request.Config) {
					if config.ScheduleTimeout() != -1*time.Second {
						t.Errorf("expected schedule timeout to be -1 * time.Second, got %s", config.ScheduleTimeout())
					}
					if config.RequestTimeout() != -1*time.Second {
						t.Errorf("expected request timeout to be -1 * time.Second, got %s", config.RequestTimeout())
					}
				},
			},
		},
		{
			name: "can override defaults",
			setup: []request.Option{
				request.WithDefaultOption(config.WithHost(otherUrl)),
			},
			expectations: []expectConfig{
				func(config *request.Config) {
					if config.Host() != otherUrl {
						t.Errorf("expected host to be https://example.com, got %s", config.Host())
					}
				},
			},
		},
		{
			name: "can set token",
			setup: []request.Option{
				request.WithToken(newMockToken(t, ctrl, "token")),
			},
			expectations: []expectConfig{
				func(config *request.Config) {
					token := config.Token()
					if token == nil {
						t.Errorf("expected token to be not nil")
					}
					if token.Token() != "token" {
						t.Errorf("expected token to be token, got %s", token.Token())
					}
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			config := request.NewConfig(nil, tc.setup...)
			for _, expectFn := range tc.expectations {
				expectFn(config)
			}
		})
	}
}

func TestCustomSender(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	otherUrl, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testDefaults := config.NewDefaults(config.WithHost(otherUrl), config.WithUserAgent("test"), config.WithCompatibilityDate("2025-01-01"))

	mockDefaultsProvider := mock.NewMockDefaultsProvider(ctrl)
	mockDefaultsProvider.EXPECT().DefaultRequestConfig().Return(testDefaults).Times(1)

	config := request.NewConfig(mockDefaultsProvider)

	if config.Host() != otherUrl {
		t.Errorf("expected host to be https://example.com, got %s", config.Host())
	}
	if config.UserAgent() != "test" {
		t.Errorf("expected user agent to be test, got %s", config.UserAgent())
	}
	if config.CompatibilityDate() != "2025-01-01" {
		t.Errorf("expected compatibility date to be 2025-01-01, got %s", config.CompatibilityDate())
	}
}

func TestApplyHeaders(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	testCases := []struct {
		name            string
		setup           []request.Option
		expectedHeaders http.Header
	}{
		{
			name: "defaults are applied",
			expectedHeaders: http.Header{
				"Accept-Language":      {defaults.Language},
				"User-Agent":           {defaults.UserAgent},
				"X-Compatibility-Date": {defaults.CompatibilityDate},
			},
		},
		{
			name: "can set headers",
			setup: []request.Option{
				request.WithHeader("X-Test", "test"),
			},
			expectedHeaders: http.Header{
				"X-Test":               {"test"},
				"Accept-Language":      {defaults.Language},
				"User-Agent":           {defaults.UserAgent},
				"X-Compatibility-Date": {defaults.CompatibilityDate},
			},
		},
		{
			name: "can set token",
			setup: []request.Option{
				request.WithToken(newMockToken(t, ctrl, "token")),
			},
			expectedHeaders: http.Header{
				"Authorization":        {"Bearer token"},
				"Accept-Language":      {defaults.Language},
				"User-Agent":           {defaults.UserAgent},
				"X-Compatibility-Date": {defaults.CompatibilityDate},
			},
		},
		{
			name: "can override defaults",
			setup: []request.Option{
				request.WithDefaultOption(config.WithUserAgent("test")),
			},
			expectedHeaders: http.Header{
				"User-Agent":           {"test"},
				"Accept-Language":      {defaults.Language},
				"X-Compatibility-Date": {defaults.CompatibilityDate},
			},
		},
		{
			name: "can set ETag",
			setup: []request.Option{
				request.WithETag("etag"),
			},
			expectedHeaders: http.Header{
				"If-None-Match":        {"etag"},
				"Accept-Language":      {defaults.Language},
				"User-Agent":           {defaults.UserAgent},
				"X-Compatibility-Date": {defaults.CompatibilityDate},
			},
		},
		{
			name: "can set LastModified",
			setup: []request.Option{
				request.WithLastModified(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedHeaders: http.Header{
				"If-Modified-Since":    {"Wed, 01 Jan 2025 00:00:00 GMT"},
				"Accept-Language":      {defaults.Language},
				"User-Agent":           {defaults.UserAgent},
				"X-Compatibility-Date": {defaults.CompatibilityDate},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			config := request.NewConfig(nil, tc.setup...)
			headers := make(http.Header, 0)
			config.ApplyHeaders(headers)
			if diff := cmp.Diff(tc.expectedHeaders, headers); diff != "" {
				t.Fatalf("headers mismatch (-want +got): %s", diff)
			}
		})
	}
}
