package faction_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/faction"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `500002`
	var v faction.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != faction.Identifier(500002) {
		t.Fatalf("got %v want %v", v, faction.Identifier(500002))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
