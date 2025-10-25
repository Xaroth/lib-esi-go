package requestgen_test

import (
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/gentest"
	"github.com/xaroth/lib-esi-go/internal/generate/requestgen"
)

func TestBuildPackage_nestedObject(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"GetUniverseStargatesStargateId"})
	if err != nil {
		t.Fatal(err)
	}
	cfg := requestgen.Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	pkg, err := requestgen.BuildPackage(ops[0], spec, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(pkg.OutputFields) < 2 {
		t.Fatalf("OutputFields: %+v", pkg.OutputFields)
	}
	found := false
	for _, f := range pkg.OutputFields {
		if f.Name == "Position" && f.Type.Type == "Position" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected Position field, got %+v", pkg.OutputFields)
	}
}
