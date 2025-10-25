package commonmodels_test

import (
	"strings"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/commonmodels"
)

func TestParsedTemplates(t *testing.T) {
	tmpl, err := commonmodels.ParsedTemplates()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"int64", "string_alias", "enum", "test"} {
		if tmpl.Lookup(name) == nil {
			t.Errorf("missing template %q", name)
		}
	}
}

func TestExecuteTemplate_int64(t *testing.T) {
	out, err := commonmodels.ExecuteTemplate("int64", commonmodels.Model{Package: "alliance"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "package alliance") {
		t.Errorf("unexpected output: %s", out)
	}
}
