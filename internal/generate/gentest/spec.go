package gentest

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

// MinimalSpecPath returns the path to the shared minimal OpenAPI fixture.
func MinimalSpecPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "openapi", "testdata", "minimal.json")
}

// LoadMinimalSpec loads the shared minimal OpenAPI fixture.
func LoadMinimalSpec(t *testing.T) *openapi.Spec {
	t.Helper()
	spec, err := openapi.LoadSpec(t.Context(), "", "", MinimalSpecPath())
	if err != nil {
		t.Fatal(err)
	}
	return spec
}
