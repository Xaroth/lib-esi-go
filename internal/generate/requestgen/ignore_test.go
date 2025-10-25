package requestgen_test

import (
	"strings"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/gentest"
	"github.com/xaroth/lib-esi-go/internal/generate/requestgen"
)

func TestBuildPackage_ignoresClientParameters(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"GetAlliancesAllianceId"})
	if err != nil {
		t.Fatal(err)
	}
	cfg := requestgen.Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	pkg, err := requestgen.BuildPackage(ops[0], spec, cfg)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range pkg.InputFields {
		if strings.Contains(f.TagVal, "Compatibility") || strings.Contains(f.TagVal, "Tenant") {
			t.Errorf("unexpected ignored field: %+v", f)
		}
	}
	if len(pkg.InputFields) != 1 || pkg.InputFields[0].TagKey != "path" {
		t.Fatalf("input fields = %+v", pkg.InputFields)
	}
}
