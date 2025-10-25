package cmdutil

import (
	"os"
	"path/filepath"
	"testing"

	defaults "github.com/xaroth/lib-esi-go"
)

func TestCompatDateFromArgs_library(t *testing.T) {
	date, rest := CompatDateFromArgs([]string{"LIBRARY", "GetFoo"}, true)
	if date != defaults.CompatibilityDate {
		t.Errorf("date = %q", date)
	}
	if len(rest) != 1 || rest[0] != "GetFoo" {
		t.Errorf("rest = %v", rest)
	}
}

func TestCompatDateFromArgs_libraryCaseInsensitive(t *testing.T) {
	date, _ := CompatDateFromArgs([]string{"library"}, true)
	if date != defaults.CompatibilityDate {
		t.Errorf("date = %q", date)
	}
}

func TestOutputDir_dotUsesCwd(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "esi")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(sub); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	out, err := OutputDir(".")
	if err != nil {
		t.Fatal(err)
	}
	if out != sub {
		t.Errorf("OutputDir(.) = %q, want %q", out, sub)
	}
}
