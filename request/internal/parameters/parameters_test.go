package parameters_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/xaroth/lib-esi-go/request/internal/parameters"
)

type queryStringer string

func (s queryStringer) String() string { return string(s) }

func TestExtract(t *testing.T) {
	t.Parallel()

	type jsonBody struct {
		K string `json:"k"`
		N int    `json:"n"`
	}

	testCases := []struct {
		name             string
		invoke           func() (map[string]any, url.Values, http.Header, io.Reader, error)
		expectedPath     map[string]any
		expectedQuery    url.Values
		expectedHeader   http.Header
		expectedBodyJSON string
		expectedErr      error
	}{
		{
			name: "all parameters with explicit json body",
			invoke: func() (map[string]any, url.Values, http.Header, io.Reader, error) {
				type input struct {
					PathID  int      `path:"id"`
					QueryA  string   `query:"a"`
					QueryB  int64    `query:"b"`
					HeaderX string   `header:"X-Test"`
					Body    jsonBody `body:"json"`
				}
				in := &input{
					PathID:  123,
					QueryA:  "foo",
					QueryB:  456,
					HeaderX: "hdr",
					Body:    jsonBody{K: "v", N: 7},
				}
				return parameters.Extract(in)
			},
			expectedPath: map[string]any{
				"id": 123,
			},
			expectedQuery: url.Values{
				"a": {"foo"},
				"b": {"456"},
			},
			expectedHeader: http.Header{
				"X-Test":       {"hdr"},
				"Content-Type": {"application/json"},
			},
			expectedBodyJSON: `{"k":"v","n":7}`,
		},
		{
			name: "json body via empty tag",
			invoke: func() (map[string]any, url.Values, http.Header, io.Reader, error) {
				type input struct {
					Body jsonBody `body:""`
				}
				in := &input{Body: jsonBody{K: "x", N: 1}}
				return parameters.Extract(in)
			},
			expectedPath:     map[string]any{},
			expectedQuery:    url.Values{},
			expectedHeader:   http.Header{"Content-Type": {"application/json"}},
			expectedBodyJSON: `{"k":"x","n":1}`,
		},
		{
			name: "duplicate query parameters",
			invoke: func() (map[string]any, url.Values, http.Header, io.Reader, error) {
				type input struct {
					Q1 string `query:"q"`
					Q2 int    `query:"q"`
				}
				in := &input{Q1: "a", Q2: 2}
				return parameters.Extract(in)
			},
			expectedPath: map[string]any{},
			expectedQuery: url.Values{
				"q": {"a", "2"},
			},
			expectedHeader:   http.Header{},
			expectedBodyJSON: "",
		},
		{
			name: "no body does not set content-type",
			invoke: func() (map[string]any, url.Values, http.Header, io.Reader, error) {
				type input struct {
					HeaderX string `header:"X-Test"`
				}
				in := &input{HeaderX: "y"}
				return parameters.Extract(in)
			},
			expectedPath:  map[string]any{},
			expectedQuery: url.Values{},
			expectedHeader: http.Header{
				"X-Test": {"y"},
			},
		},
		{
			name: "path accepts any type",
			invoke: func() (map[string]any, url.Values, http.Header, io.Reader, error) {
				type input struct {
					V any `path:"v"`
				}
				in := &input{V: 3.14}
				return parameters.Extract(in)
			},
			expectedPath: map[string]any{
				"v": 3.14,
			},
			expectedQuery:    url.Values{},
			expectedHeader:   http.Header{},
			expectedBodyJSON: "",
		},
		{
			name: "query fmt.Stringer value",
			invoke: func() (map[string]any, url.Values, http.Header, io.Reader, error) {
				type input struct {
					Q queryStringer `query:"q"`
				}
				in := &input{Q: "hello"}
				return parameters.Extract(in)
			},
			expectedPath:  map[string]any{},
			expectedQuery: url.Values{"q": {"hello"}},
			expectedHeader: http.Header{},
		},
		{
			name: "header fmt.Stringer value",
			invoke: func() (map[string]any, url.Values, http.Header, io.Reader, error) {
				type input struct {
					H queryStringer `header:"X-H"`
				}
				in := &input{H: "hdr-val"}
				return parameters.Extract(in)
			},
			expectedPath: map[string]any{},
			expectedQuery: url.Values{},
			expectedHeader: http.Header{
				"X-H": {"hdr-val"},
			},
		},
		{
			name: "invalid query value type",
			invoke: func() (map[string]any, url.Values, http.Header, io.Reader, error) {
				type input struct {
					Q bool `query:"q"`
				}
				in := &input{Q: true}
				return parameters.Extract(in)
			},
			expectedErr: parameters.ErrInvalidValueType,
		},
		{
			name: "invalid header value type",
			invoke: func() (map[string]any, url.Values, http.Header, io.Reader, error) {
				type input struct {
					H bool `header:"X-H"`
				}
				in := &input{H: false}
				return parameters.Extract(in)
			},
			expectedErr: parameters.ErrInvalidValueType,
		},
		{
			name: "invalid body type tag",
			invoke: func() (map[string]any, url.Values, http.Header, io.Reader, error) {
				type input struct {
					Body jsonBody `body:"xml"`
				}
				in := &input{Body: jsonBody{K: "z", N: 9}}
				return parameters.Extract(in)
			},
			expectedErr: parameters.ErrInvalidBodyType,
		},
		{
			name: "json marshal failure in body",
			invoke: func() (map[string]any, url.Values, http.Header, io.Reader, error) {
				type bad struct {
					Ch chan int `json:"ch"`
				}
				type input struct {
					Body any `body:"json"`
				}
				in := &input{Body: bad{Ch: make(chan int)}}
				return parameters.Extract(in)
			},
			// Expect a JSON marshaling error; we'll just assert non-nil.
			expectedErr: errors.New("json error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			path, query, header, body, err := tc.invoke()

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				// If expected is a known sentinel, use errors.Is
				if errors.Is(tc.expectedErr, parameters.ErrInvalidValueType) || errors.Is(tc.expectedErr, parameters.ErrInvalidBodyType) {
					if !errors.Is(err, tc.expectedErr) {
						t.Fatalf("got error %v, want %v", err, tc.expectedErr)
					}
				}
				// For other error types (e.g., JSON marshal), just assert non-nil.
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.expectedPath, path); diff != "" {
				t.Fatalf("path parameters mismatch (-want +got): %s", diff)
			}
			if diff := cmp.Diff(tc.expectedQuery, query); diff != "" {
				t.Fatalf("query parameters mismatch (-want +got): %s", diff)
			}
			if diff := cmp.Diff(tc.expectedHeader, header); diff != "" {
				t.Fatalf("header parameters mismatch (-want +got): %s", diff)
			}

			switch {
			case tc.expectedBodyJSON == "" && body != nil:
				t.Fatalf("expected no body, got one")
			case tc.expectedBodyJSON != "":
				if body == nil {
					t.Fatalf("expected body, got nil")
				}
				data, readErr := io.ReadAll(body)
				if readErr != nil {
					t.Fatalf("failed to read body: %v", readErr)
				}
				if !bytes.Equal([]byte(tc.expectedBodyJSON), data) {
					t.Fatalf("body mismatch: got %q, want %q", string(data), tc.expectedBodyJSON)
				}
			}
		})
	}
}
