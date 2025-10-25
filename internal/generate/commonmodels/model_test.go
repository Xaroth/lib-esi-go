package commonmodels_test

import (
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/commonmodels"
	"github.com/xaroth/lib-esi-go/internal/generate/gentest"
)

func TestModelsFromSpec(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	models, err := commonmodels.ModelsFromSpec(spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(models) == 0 {
		t.Fatal("expected common models")
	}
	found := false
	for _, m := range models {
		if m.SchemaName == "AllianceID" {
			found = true
			if m.Package != "alliance" || m.TypeName != "Identifier" {
				t.Errorf("AllianceID model = %+v", m)
			}
		}
	}
	if !found {
		t.Fatalf("models = %+v", models)
	}
}
