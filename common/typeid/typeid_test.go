package typeid_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/typeid"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `587`
	var v typeid.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != typeid.Identifier(587) {
		t.Fatalf("got %v want %v", v, typeid.Identifier(587))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
