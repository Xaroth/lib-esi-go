package requestgen_test

import (
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/gentest"
	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
	"github.com/xaroth/lib-esi-go/internal/generate/requestgen"
)

func TestTypeMapper_commonModel(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	mapper := requestgen.NewTypeMapper(requestgen.Config{LibModule: "github.com/example/lib", CommonSuffix: "common"}, spec)
	gt, name, err := mapper.MapSchemaRef(openapi.SchemaRef{Ref: "#/components/schemas/AllianceID"}, true)
	if err != nil {
		t.Fatal(err)
	}
	if name != "AllianceID" || gt.Type != "alliance.Identifier" {
		t.Errorf("got %q %q", gt.Type, name)
	}
	if gt.Import != "github.com/example/lib/common/alliance" {
		t.Errorf("import %q", gt.Import)
	}
}

func TestTypeMapper_optionalPointer(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	mapper := requestgen.NewTypeMapper(requestgen.Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}, spec)
	gt, _, err := mapper.MapSchemaRef(openapi.SchemaRef{Ref: "#/components/schemas/FactionID"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if gt.Type != "*faction.Identifier" {
		t.Errorf("got %q", gt.Type)
	}
}
