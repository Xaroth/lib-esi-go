package requestgen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/gentest"
	"github.com/xaroth/lib-esi-go/internal/generate/requestgen"
)

func TestWritePackages(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"GetUniverseFactions"})
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	cfg := requestgen.Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	pkg, err := requestgen.BuildPackage(ops[0], spec, cfg)
	if err != nil {
		t.Fatal(err)
	}
	n, err := requestgen.WritePackages(dir, []requestgen.PackageModel{pkg}, cfg, false)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Errorf("written = %d", n)
	}
	if _, err := os.Stat(filepath.Join(dir, "getuniversefactions", "request.go")); err != nil {
		t.Fatal(err)
	}
}

func TestWritePackages_check(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"GetUniverseFactions"})
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	cfg := requestgen.Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	if _, err := requestgen.BuildAndWrite(spec, ops, dir, cfg, false); err != nil {
		t.Fatal(err)
	}
	if _, err := requestgen.BuildAndWrite(spec, ops, dir, cfg, true); err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(filepath.Join(dir, "getuniversefactions", "request.go"))
	if _, err := requestgen.BuildAndWrite(spec, ops, dir, cfg, true); err == nil {
		t.Fatal("expected check error for missing file")
	}
}
