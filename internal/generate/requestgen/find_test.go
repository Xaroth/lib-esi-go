package requestgen_test

import (
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/gentest"
	"github.com/xaroth/lib-esi-go/internal/generate/requestgen"
)

func TestFindOperations_byOperationID(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"GetAlliancesAllianceId"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 || ops[0].Path != "/alliances/{alliance_id}" || ops[0].Method != "GET" {
		t.Fatalf("got %+v", ops[0])
	}
}

func TestFindOperations_byMethodPath(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"get /universe/factions"})
	if err != nil {
		t.Fatal(err)
	}
	if ops[0].OperationID != "GetUniverseFactions" {
		t.Fatalf("got %+v", ops[0])
	}
}

func TestFindOperations_caseInsensitive(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	_, err := requestgen.FindOperations(spec, []string{"getalliancesallianceid"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindOperations_notFound(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	_, err := requestgen.FindOperations(spec, []string{"NoSuchOp"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFindOperations_allPaths(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	ops, err := requestgen.FindOperations(spec, []string{"ALL_PATHS"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 7 {
		t.Fatalf("got %d operations, want 7", len(ops))
	}
}

func TestFindOperations_allPathsNotExclusive(t *testing.T) {
	spec := gentest.LoadMinimalSpec(t)
	_, err := requestgen.FindOperations(spec, []string{"ALL_PATHS", "GetUniverseFactions"})
	if err == nil {
		t.Fatal("expected error when ALL_PATHS is combined with other selectors")
	}
}
