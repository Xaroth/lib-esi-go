package shiptreegroup_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/shiptreegroup"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `1559`
	var v shiptreegroup.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != shiptreegroup.Identifier(1559) {
		t.Fatalf("got %v want %v", v, shiptreegroup.Identifier(1559))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
