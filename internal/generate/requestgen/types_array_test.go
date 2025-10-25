package requestgen_test

import (
	"strings"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/gentest"
	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
	"github.com/xaroth/lib-esi-go/internal/generate/requestgen"
)

func TestMapSchema_arrayCommonModelPreservesImport(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	cfg := requestgen.Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	mapper := requestgen.NewTypeMapper(cfg, spec)

	gt, _, err := mapper.MapSchemaRef(openapi.SchemaRef{
		Type: "array",
		Items: &openapi.SchemaRef{
			Ref: "#/components/schemas/AllianceID",
		},
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	if gt.Type != "[]alliance.Identifier" {
		t.Fatalf("type %q", gt.Type)
	}
	if gt.Import != "github.com/xaroth/lib-esi-go/common/alliance" {
		t.Fatalf("import %q", gt.Import)
	}
}

func TestGeneratePackage_arrayCommonModelImport(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	spec.Paths["/meta/compatibility-dates"] = openapi.PathItem{
		"get": {
			OperationID: "GetMetaCompatibilityDates",
			Responses: map[string]openapi.Response{
				"200": {
					Content: map[string]openapi.MediaTypeObject{
						"application/json": {
							Schema: openapi.SchemaRef{
								Type: "object",
								Properties: map[string]openapi.SchemaRef{
									"compatibility_dates": {
										Type: "array",
										Items: &openapi.SchemaRef{
											Ref: "#/components/schemas/CompatibilityDate",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	spec.Components.Schemas["CompatibilityDate"] = openapi.Schema{
		XCommonModel: openapi.Boolish(true),
		Type:       "string",
		Format:     "date",
	}

	ops, err := requestgen.FindOperations(spec, []string{"GetMetaCompatibilityDates"})
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
	if !strings.Contains(output, "common/compatibilitydate") {
		t.Errorf("output missing import: %s", output)
	}
}
