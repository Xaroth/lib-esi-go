package requestgen_test

import (
	"strings"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/gentest"
	"github.com/xaroth/lib-esi-go/internal/generate/requestgen"
)

func TestGeneratePackage_oneOf(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"GetCorporationsProjectsDetail"})
	if err != nil {
		t.Fatal(err)
	}
	cfg := requestgen.Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	pkg, err := requestgen.BuildPackage(ops[0], spec, cfg)
	if err != nil {
		t.Fatal(err)
	}
	files, err := requestgen.GeneratePackage(pkg, cfg)
	if err != nil {
		t.Fatal(err)
	}
	output := string(files.Output)
	for _, want := range []string{
		"type Configuration struct",
		"CorporationsProjectsDetailConfigurationmanual",
		"CorporationsProjectsDetailConfigurationdamageship",
		`json:"manual,omitempty"`,
		`json:"damage_ship,omitempty"`,
		"type Identity struct",
		"character.Identifier",
		"corporation.Identifier",
		"Identities []Identity",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}
