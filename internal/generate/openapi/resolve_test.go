package openapi_test

import (
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

func TestResolveParameterRef(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Parameters: map[string]openapi.Parameter{
				"AllianceId": {
					Name:     "alliance_id",
					In:       "path",
					Required: true,
					Schema:   &openapi.SchemaRef{Ref: "#/components/schemas/AllianceID"},
				},
			},
			Schemas: map[string]openapi.Schema{
				"AllianceID": {Type: "integer", Format: "int64", XCommonModel: openapi.Boolish(true)},
			},
		},
	}
	r := openapi.NewResolver(spec)

	p, name, err := r.ResolveParameterRef(openapi.ParameterRef{Ref: "#/components/parameters/AllianceId"})
	if err != nil {
		t.Fatal(err)
	}
	if name != "AllianceId" {
		t.Errorf("name = %q", name)
	}
	if p.Name != "alliance_id" || p.In != "path" {
		t.Errorf("parameter = %+v", p)
	}
}

func TestResolveSchemaRef(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Schemas: map[string]openapi.Schema{
				"AllianceID": {Type: "integer", Format: "int64", XCommonModel: openapi.Boolish(true)},
			},
		},
	}
	r := openapi.NewResolver(spec)

	s, name, err := r.ResolveSchemaRef(openapi.SchemaRef{Ref: "#/components/schemas/AllianceID"})
	if err != nil {
		t.Fatal(err)
	}
	if name != "AllianceID" {
		t.Errorf("name = %q", name)
	}
	if s.Type != "integer" || !s.XCommonModel.IsTrue() {
		t.Errorf("schema = %+v", s)
	}
}

func TestResolveSchemaRef_inline(t *testing.T) {
	spec := &openapi.Spec{}
	r := openapi.NewResolver(spec)

	s, name, err := r.ResolveSchemaRef(openapi.SchemaRef{Type: "string", Format: "uuid"})
	if err != nil {
		t.Fatal(err)
	}
	if name != "" || s.Type != "string" {
		t.Errorf("got name=%q schema=%+v", name, s)
	}
}

func TestResolveSchemaRef_missing(t *testing.T) {
	spec := &openapi.Spec{Components: openapi.Components{Schemas: map[string]openapi.Schema{}}}
	r := openapi.NewResolver(spec)

	_, _, err := r.ResolveSchemaRef(openapi.SchemaRef{Ref: "#/components/schemas/Missing"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCommonSchemaNames(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Schemas: map[string]openapi.Schema{
				"AllianceID":  {XCommonModel: openapi.Boolish(true)},
				"OtherSchema": {XCommonModel: openapi.Boolish(false)},
				"FactionID":   {XCommonModel: openapi.Boolish(true)},
			},
		},
	}
	names := openapi.CommonSchemaNames(spec)
	if len(names) != 2 {
		t.Fatalf("got %v", names)
	}
}
