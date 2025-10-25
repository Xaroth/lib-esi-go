package request_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/xaroth/lib-esi-go/request"
)

type testPatternVar struct{}

func (testPatternVar) PatternVariable() string {
	return "pv"
}

func TestPatternString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		path        string
		variables   map[string]any
		expected    string
		expectedErr error
	}{
		{
			name:      "literal only",
			path:      "a/b/c",
			variables: map[string]any{},
			expected:  "a/b/c",
		},
		{
			name: "string variable",
			path: "a/{x}/c",
			variables: map[string]any{
				"x": "Y",
			},
			expected: "a/Y/c",
		},
		{
			name: "int variable",
			path: "{id}",
			variables: map[string]any{
				"id": 42,
			},
			expected: "42",
		},
		{
			name: "int64 variable",
			path: "p/{big}",
			variables: map[string]any{
				"big": int64(1234567890123),
			},
			expected: "p/1234567890123",
		},
		{
			name: "float64 variable",
			path: "a/{x}",
			variables: map[string]any{
				"x": 3.14,
			},
			expected: "a/3.14",
		},
		{
			name: "float without trailing zeros",
			path: "p/{x}",
			variables: map[string]any{
				"x": 3.0,
			},
			expected: "p/3",
		},
		{
			name: "pattern variable interface",
			path: "x/{pv}/y",
			variables: map[string]any{
				"pv": testPatternVar{},
			},
			expected: "x/pv/y",
		},
		{
			name:      "missing variable",
			path:      "a/{x}/b",
			variables: map[string]any{
				// missing "x"
			},
			expectedErr: request.ErrMissingVariable,
		},
		{
			name: "invalid type",
			path: "a/{x}",
			variables: map[string]any{
				"x": true, // unsupported type
			},
			expectedErr: request.ErrInvalidVariable,
		},
		{
			name: "extraneous variables",
			path: "a/{x}",
			variables: map[string]any{
				"x": "ok",
				"y": "extra",
			},
			expectedErr: request.ErrExtraneousVariable,
		},
		{
			name: "suffix slash preserved",
			path: "a/{x}/",
			variables: map[string]any{
				"x": "t",
			},
			expected: "a/t/",
		},
		{
			name: "prefix slash preserved",
			path: "/a/{x}",
			variables: map[string]any{
				"x": "t",
			},
			expected: "/a/t",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := request.NewPattern(tc.path)
			got, err := p.String(tc.variables)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.expectedErr)
				}
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("error %v is not %v", err, tc.expectedErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Fatalf("result mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestPatternVariables(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		path              string
		expectedVariables []string
	}{
		{
			name:              "literal only",
			path:              "a/b/c",
			expectedVariables: []string{},
		},
		{
			name:              "single variable",
			path:              "a/{x}/c",
			expectedVariables: []string{"x"},
		},
		{
			name:              "multiple variables",
			path:              "a/{x}/b/{y}",
			expectedVariables: []string{"x", "y"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := request.NewPattern(tc.path)
			got := p.Variables()
			if diff := cmp.Diff(tc.expectedVariables, got); diff != "" {
				t.Fatalf("result mismatch (-want +got): %s", diff)
			}
		})
	}
}
