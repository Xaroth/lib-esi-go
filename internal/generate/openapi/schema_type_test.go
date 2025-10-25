package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

func TestSchemaType_UnmarshalJSON_string(t *testing.T) {
	var ref openapi.SchemaRef
	if err := json.Unmarshal([]byte(`{"type":"string"}`), &ref); err != nil {
		t.Fatal(err)
	}
	if ref.Type != "string" {
		t.Fatalf("got %q", ref.Type)
	}
}

func TestSchemaType_UnmarshalJSON_nullableArray(t *testing.T) {
	var ref openapi.SchemaRef
	if err := json.Unmarshal([]byte(`{"type":["boolean","null"]}`), &ref); err != nil {
		t.Fatal(err)
	}
	if ref.Type != "boolean" {
		t.Fatalf("got %q", ref.Type)
	}
}
