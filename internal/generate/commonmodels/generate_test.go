package commonmodels_test

import (
	"strings"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/commonmodels"
	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

func TestGeneratePackage_int64ID(t *testing.T) {
	m := commonmodels.Model{
		SchemaName: "AllianceID",
		Package:    "alliance",
		TypeName:   "Identifier",
		Schema: openapi.Schema{
			Type:   "integer",
			Format: "int64",
		},
		Examples: []any{float64(99000001)},
	}
	mainGo, testGo, err := commonmodels.GeneratePackage(m, "github.com/xaroth/lib-esi-go", "common")
	if err != nil {
		t.Fatal(err)
	}
	main := string(mainGo)
	if !strings.Contains(main, "type Identifier int64") {
		t.Errorf("main missing Identifier: %s", main)
	}
	if !strings.Contains(main, "func (id Identifier) String() string") {
		t.Errorf("main missing String(): %s", main)
	}
	if testGo == nil {
		t.Fatal("expected test file")
	}
	if !strings.Contains(string(testGo), "JSONRoundTrip") {
		t.Errorf("test missing round trip: %s", testGo)
	}
	if !strings.Contains(string(testGo), "package alliance_test") {
		t.Errorf("test should use external package alliance_test: %s", testGo)
	}
	if !strings.Contains(string(testGo), "alliance.Identifier") {
		t.Errorf("test should qualify types: %s", testGo)
	}
	if !strings.Contains(string(testGo), "github.com/xaroth/lib-esi-go/common/alliance") {
		t.Errorf("test should import module path: %s", testGo)
	}
}

func TestGeneratePackage_uuid(t *testing.T) {
	m := commonmodels.Model{
		SchemaName: "UUID",
		Package:    "uuid",
		TypeName:   "UUID",
		Schema: openapi.Schema{
			Type:   "string",
			Format: "uuid",
		},
		Examples: []any{"3868eaed-8278-4cb7-9709-7d7de9c20dc7"},
	}
	mainGo, _, err := commonmodels.GeneratePackage(m, "github.com/xaroth/lib-esi-go", "common")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(mainGo), "type UUID string") {
		t.Errorf("unexpected main: %s", mainGo)
	}
	if !strings.Contains(string(mainGo), "String() string") {
		t.Errorf("main missing String(): %s", mainGo)
	}
}
