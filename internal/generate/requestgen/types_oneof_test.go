package requestgen

import (
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/gentest"
	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

func TestMapSchema_untypedJSON(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	mapper := NewTypeMapper(Config{LibModule: "github.com/example/lib", CommonSuffix: "common"}, spec)
	gt, _, err := mapper.MapSchema(openapi.Schema{}, "", false)
	if err != nil {
		t.Fatal(err)
	}
	if gt.Type != "json.RawMessage" {
		t.Fatalf("got %q", gt.Type)
	}
}
