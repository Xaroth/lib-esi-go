package commonmodels_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/commonmodels"
	"github.com/xaroth/lib-esi-go/internal/generate/gentest"
)

func TestWritePackages(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	models, err := commonmodels.ModelsFromSpec(spec)
	if err != nil {
		t.Fatal(err)
	}
	var alliance commonmodels.Model
	for _, m := range models {
		if m.SchemaName == "AllianceID" {
			alliance = m
			break
		}
	}
	dir := moduleTempDir(t)
	n, err := commonmodels.WritePackages(dir, []commonmodels.Model{alliance}, false)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("written = %d", n)
	}
	if _, err := os.Stat(filepath.Join(dir, "alliance", "alliance.go")); err != nil {
		t.Fatal(err)
	}
}

func TestWritePackages_check(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	models, err := commonmodels.ModelsFromSpec(spec)
	if err != nil {
		t.Fatal(err)
	}
	dir := moduleTempDir(t)
	if _, err := commonmodels.WritePackages(dir, models, false); err != nil {
		t.Fatal(err)
	}
	if _, err := commonmodels.WritePackages(dir, models, true); err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(filepath.Join(dir, "alliance", "alliance.go"))
	if _, err := commonmodels.WritePackages(dir, models, true); err == nil {
		t.Fatal("expected check error for missing file")
	}
}

func moduleTempDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module github.com/xaroth/lib-esi-go\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestModelsFromSpec_nil(t *testing.T) {
	if _, err := commonmodels.ModelsFromSpec(nil); err == nil {
		t.Fatal("expected error")
	}
}
