package requestgen

import "testing"

func TestSingularStructName(t *testing.T) {
	tests := []struct {
		wire string
		want string
	}{
		{"identities", "Identity"},
		{"locations", "Location"},
		{"configuration", "Configuration"},
		{"docking_locations", "DockingLocation"},
	}
	for _, tc := range tests {
		if got := singularStructName(tc.wire); got != tc.want {
			t.Errorf("singularStructName(%q) = %q, want %q", tc.wire, got, tc.want)
		}
	}
}
