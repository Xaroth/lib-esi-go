package request_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	"github.com/xaroth/lib-esi-go/internal/config"
	"github.com/xaroth/lib-esi-go/request"
	"github.com/xaroth/lib-esi-go/request/mock"
	"github.com/xaroth/lib-esi-go/util/esierror"
)

//go:generate go run -mod=mod go.uber.org/mock/mockgen -build_flags=--mod=mod -destination mock/mock_sender.go -package mock github.com/xaroth/lib-esi-go/request SchedulingSender
//go:generate go run -mod=mod go.uber.org/mock/mockgen -build_flags=--mod=mod -destination mock/mock_token.go -package mock github.com/xaroth/lib-esi-go/request Token,RefreshableToken

type bodyStruct struct {
	K string `json:"k"`
	N int    `json:"n"`
}

type requestInput struct {
	PathID  int        `path:"id"`
	QueryA  string     `query:"a"`
	QueryB  int64      `query:"b"`
	HeaderX string     `header:"X-Test"`
	Body    bodyStruct `body:"json"`
}

type requestOutput struct {
	ID int    `json:"id"`
	A  string `json:"a"`
	B  int64  `json:"b"`
	X  string `json:"x"`
	K  string `json:"k"`
	N  int    `json:"n"`
}

type expectation func(r *http.Request) error

func defaultResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id":123,"a":"foo","b":456,"x":"hdr","k":"v","n":7}`))
}

var defaultRequestInput = &requestInput{
	PathID:  123,
	QueryA:  "foo",
	QueryB:  456,
	HeaderX: "hdr",
	Body: bodyStruct{
		K: "v",
		N: 7,
	},
}

var defaultResponseOutput = &requestOutput{
	ID: 123,
	A:  "foo",
	B:  456,
	X:  "hdr",
	K:  "v",
	N:  7,
}

type client struct {
	*mock.MockSchedulingSender
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

var defaultExpectations = []expectation{
	func(r *http.Request) error {
		if r.Method != "GET" {
			return fmt.Errorf("unexpected method: %s", r.Method)
		}
		return nil
	},
	func(r *http.Request) error {
		if r.URL.Path != "/test/123" {
			return fmt.Errorf("unexpected path: %s", r.URL.Path)
		}
		return nil
	},
	func(r *http.Request) error {
		if r.URL.Query().Get("a") != "foo" {
			return fmt.Errorf("unexpected query: %s", r.URL.Query().Get("a"))
		}
		return nil
	},
	func(r *http.Request) error {
		if r.Header.Get("X-Test") != "hdr" {
			return fmt.Errorf("unexpected header: %s", r.Header.Get("X-Test"))
		}
		return nil
	},
}

func TestCreate(t *testing.T) {
	t.Parallel()

	builder := request.Create[requestInput, requestOutput](request.RequestInfo{
		Method: "GET",
		Path:   "/test/{id}",
	})

	testCases := []struct {
		name             string
		input            *requestInput
		handler          func(w http.ResponseWriter, r *http.Request)
		tokenFn          func(ctrl *gomock.Controller) request.Token
		expectations     []expectation
		expectedResponse *requestOutput
		expectedErr      error
	}{
		{
			name:             "simple request",
			input:            defaultRequestInput,
			tokenFn:          func(ctrl *gomock.Controller) request.Token { return nil },
			expectations:     defaultExpectations,
			expectedResponse: defaultResponseOutput,
		},
		{
			name:    "error response returns ErrInvalidStatusCode",
			input:   defaultRequestInput,
			tokenFn: func(ctrl *gomock.Controller) request.Token { return nil },
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"internal server error"}`))
			},
			expectations:     []expectation{},
			expectedErr:      request.ErrInvalidStatusCode,
			expectedResponse: nil,
		},
		{
			name:    "error response returns ErrorData object",
			input:   defaultRequestInput,
			tokenFn: func(ctrl *gomock.Controller) request.Token { return nil },
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"internal server error", "details": [{"message": "internal server error", "location": "server", "value": "server"}]}`))
			},
			expectations: []expectation{},
			expectedErr: &esierror.ErrorData{
				ErrorMessage: "internal server error",
				Details: []esierror.ErrorDetails{
					{
						Message:  "internal server error",
						Location: "server",
						Value:    "server",
					},
				},
			},
			expectedResponse: nil,
		},
		{
			name:  "tokens should be passed to the request",
			input: defaultRequestInput,
			tokenFn: func(ctrl *gomock.Controller) request.Token {
				token := mock.NewMockToken(ctrl)
				token.EXPECT().Token().Return("token").Times(1)
				token.EXPECT().Owner().Return(int64(1)).AnyTimes()
				return token
			},
			expectations: []expectation{
				func(r *http.Request) error {
					if r.Header.Get("Authorization") != "Bearer token" {
						return fmt.Errorf("unexpected authorization header: %s", r.Header.Get("Authorization"))
					}
					return nil
				},
			},
			expectedResponse: defaultResponseOutput,
		},
		{
			name:  "refreshable token should be refreshed if needed",
			input: defaultRequestInput,
			tokenFn: func(ctrl *gomock.Controller) request.Token {
				token := mock.NewMockRefreshableToken(ctrl)
				token.EXPECT().Token().Return("token").Times(1)
				token.EXPECT().Owner().Return(int64(1)).AnyTimes()
				token.EXPECT().RefreshIfNeeded(gomock.Any()).Return(nil).Times(1)
				return token
			},
			expectations: []expectation{
				func(r *http.Request) error {
					if r.Header.Get("Authorization") != "Bearer token" {
						return fmt.Errorf("unexpected authorization header: %s", r.Header.Get("Authorization"))
					}
					return nil
				},
			},
			expectedResponse: defaultResponseOutput,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			token := tc.tokenFn(ctrl)

			ctx := t.Context()

			var notified bool = false
			notifyDone := func(resp *http.Response) {
				notified = true
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for _, expectation := range tc.expectations {
					if err := expectation(r); err != nil {
						t.Fatalf("expectation failed: %v", err)
					}
				}
				if tc.handler != nil {
					tc.handler(w, r)
				} else {
					defaultResponse(w, r)
				}
			}))
			defer server.Close()

			mockClient := mock.NewMockSchedulingSender(ctrl)
			mockClient.EXPECT().Schedule(gomock.Any(), gomock.Any(), gomock.Any()).Return(notifyDone, nil).Times(1)
			client := &client{MockSchedulingSender: mockClient}

			url, err := url.Parse(server.URL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			options := []request.Option{
				request.WithDefaultOption(config.WithHost(url)),
			}
			if token != nil {
				options = append(options, request.WithToken(token))
			}

			resp, err := builder(ctx, client, tc.input, options...)

			if tc.expectedErr != nil {
				if expected, ok := tc.expectedErr.(*esierror.ErrorData); ok {
					if diff := cmp.Diff(expected, resp.ErrorData); diff != "" {
						t.Fatalf("error data mismatch (-want +got): %s", diff)
					}
				} else if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.ErrorData != nil {
				t.Fatalf("expected no error data, got: %+v", resp.ErrorData)
			}

			if diff := cmp.Diff(tc.expectedResponse, resp.Data); diff != "" {
				t.Fatalf("response data mismatch (-want +got): %s", diff)
			}

			if !notified {
				t.Fatalf("expected notification to be called, but it was not")
			}
		})
	}
}
