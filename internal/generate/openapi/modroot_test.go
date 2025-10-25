package openapi_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

func TestModulePath(t *testing.T) {
	path, err := openapi.ModulePath(".")
	if err != nil {
		t.Fatal(err)
	}
	if path != "github.com/xaroth/lib-esi-go" {
		t.Errorf("got %q", path)
	}
}

func TestModulePathFromFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/foo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	path, err := openapi.ModulePathFromFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if path != "example.com/foo" {
		t.Errorf("got %q", path)
	}
}

func TestImportBase(t *testing.T) {
	root, err := openapi.ModuleRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	base, err := openapi.ImportBase(root, filepath.Join(root, "common"))
	if err != nil {
		t.Fatal(err)
	}
	if base != "common" {
		t.Errorf("ImportBase() = %q, want common", base)
	}
}
