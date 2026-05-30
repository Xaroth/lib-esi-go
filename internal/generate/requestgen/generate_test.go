package requestgen_test

import (
	"strings"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/gentest"
	"github.com/xaroth/lib-esi-go/internal/generate/requestgen"
)

func TestGeneratePackage_alliance(t *testing.T) {
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
	files, err := requestgen.GeneratePackage(pkg, cfg)
	if err != nil {
		t.Fatal(err)
	}

	input := string(files.Input)
	if !strings.Contains(input, "Alliance alliance.Identifier") || !strings.Contains(input, `path:"alliance_id"`) {
		t.Errorf("input: %s", input)
	}
	if !strings.Contains(input, "common/alliance") {
		t.Errorf("input import: %s", input)
	}
	if strings.Contains(input, `"time"`) {
		t.Errorf("input should not import time: %s", input)
	}

	output := string(files.Output)
	if !strings.Contains(output, "CreatorCorporation") || !strings.Contains(output, "corporation.Identifier") {
		t.Errorf("output: %s", output)
	}
	if !strings.Contains(output, "time.Time") || !strings.Contains(output, "DateFounded") {
		t.Errorf("output missing time: %s", output)
	}

	req := string(files.Request)
	if !strings.Contains(req, "var Request = request.Create[Input, *Output]") {
		t.Errorf("request: %s", req)
	}
	if !strings.Contains(req, `"/alliances/{alliance_id}"`) {
		t.Errorf("request path: %s", req)
	}
}

func TestGeneratePackage_arrayIDs(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"GetAlliances"})
	if err != nil {
		t.Fatal(err)
	}
	cfg := requestgen.Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	pkg, err := requestgen.BuildPackage(ops[0], spec, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if pkg.NoOutputFile || pkg.OutputType != "Output" || pkg.OutputAlias != "[]alliance.Identifier" {
		t.Fatalf("NoOutputFile=%v OutputType=%q OutputAlias=%q", pkg.NoOutputFile, pkg.OutputType, pkg.OutputAlias)
	}
	files, err := requestgen.GeneratePackage(pkg, cfg)
	if err != nil {
		t.Fatal(err)
	}
	output := string(files.Output)
	if !strings.Contains(output, "type Output = []alliance.Identifier") {
		t.Errorf("output: %s", output)
	}
	if !strings.Contains(output, "common/alliance") {
		t.Errorf("output imports: %s", output)
	}
	req := string(files.Request)
	if !strings.Contains(req, "CreateStatic[Output]") {
		t.Errorf("request: %s", req)
	}
}

func TestGeneratePackage_staticFactions(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"GetUniverseFactions"})
	if err != nil {
		t.Fatal(err)
	}
	cfg := requestgen.Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	pkg, err := requestgen.BuildPackage(ops[0], spec, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !pkg.Static {
		t.Fatal("expected static")
	}
	files, err := requestgen.GeneratePackage(pkg, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if files.Input != nil {
		t.Fatal("expected no input.go")
	}
	if !strings.Contains(string(files.Request), "CreateStatic[[]*Output]") {
		t.Errorf("request: %s", files.Request)
	}
}

func TestGeneratePackage_postAffiliation(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"PostCharactersAffiliation"})
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
	if !strings.Contains(string(files.Input), "[]int64") && !strings.Contains(string(files.Input), "body:") {
		t.Errorf("input: %s", files.Input)
	}
	if !strings.Contains(string(files.Request), "[]*Output") {
		t.Errorf("request: %s", files.Request)
	}
}

func TestGeneratePackage_nestedObject(t *testing.T) {
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
	if len(pkg.OutputNested) != 1 || pkg.OutputNested[0].Name != "Position" {
		t.Fatalf("OutputNested: %+v", pkg.OutputNested)
	}
	files, err := requestgen.GeneratePackage(pkg, cfg)
	if err != nil {
		t.Fatal(err)
	}
	output := string(files.Output)
	if !strings.Contains(output, "Position Position") || !strings.Contains(output, "type Position struct") {
		t.Errorf("output: %s", output)
	}
	if !strings.Contains(output, "X float64") || !strings.Contains(output, `json:"x"`) {
		t.Errorf("output missing position fields: %s", output)
	}
}

func TestGeneratePackage_noContent204(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"DeleteCharactersCharacterIdFittingsFittingId"})
	if err != nil {
		t.Fatal(err)
	}
	cfg := requestgen.Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	pkg, err := requestgen.BuildPackage(ops[0], spec, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !pkg.NoOutputFile || pkg.OutputType != "struct{}" {
		t.Fatalf("pkg: NoOutputFile=%v OutputType=%q", pkg.NoOutputFile, pkg.OutputType)
	}
	files, err := requestgen.GeneratePackage(pkg, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if files.Output != nil {
		t.Fatal("expected no output.go")
	}
	req := string(files.Request)
	if !strings.Contains(req, "request.Create[Input, struct{}]") {
		t.Errorf("request: %s", req)
	}
	if !strings.Contains(req, `"/characters/{character_id}/fittings/{fitting_id}"`) {
		t.Errorf("request path: %s", req)
	}
	if !strings.Contains(req, `request.WithRequiredScope("esi-fittings.write_fittings.v1")`) {
		t.Errorf("request missing required scope: %s", req)
	}
}

func TestGeneratePackage_noRequiredScope(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"GetAlliances"})
	if err != nil {
		t.Fatal(err)
	}
	cfg := requestgen.Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	pkg, err := requestgen.BuildPackage(ops[0], spec, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(pkg.RequiredScopes) != 0 {
		t.Fatalf("RequiredScopes = %v, want none", pkg.RequiredScopes)
	}
	files, err := requestgen.GeneratePackage(pkg, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(files.Request), "WithRequiredScope") {
		t.Errorf("public request should not declare scopes: %s", files.Request)
	}
}
