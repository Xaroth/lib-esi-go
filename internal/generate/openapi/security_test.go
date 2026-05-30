package openapi_test

import (
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

func TestRequiredOAuth2Scopes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		security []openapi.SecurityRequirement
		want     []string
	}{
		{
			name: "single requirement",
			security: []openapi.SecurityRequirement{
				{"OAuth2": {"esi-fittings.read_fittings.v1"}},
			},
			want: []string{"esi-fittings.read_fittings.v1"},
		},
		{
			name: "multiple scopes in one requirement",
			security: []openapi.SecurityRequirement{
				{"OAuth2": {"scope-a", "scope-b"}},
			},
			want: []string{"scope-a", "scope-b"},
		},
		{
			name: "deduplicates across requirements",
			security: []openapi.SecurityRequirement{
				{"OAuth2": {"scope-a"}},
				{"OAuth2": {"scope-b", "scope-a"}},
			},
			want: []string{"scope-a", "scope-b"},
		},
		{
			name:     "empty",
			security: nil,
			want:     nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := openapi.RequiredOAuth2Scopes(tc.security)
			if len(got) != len(tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("got %v, want %v", got, tc.want)
				}
			}
		})
	}
}
