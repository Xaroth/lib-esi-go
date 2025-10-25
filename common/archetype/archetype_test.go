package archetype_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/archetype"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `33`
	var v archetype.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != archetype.Identifier(33) {
		t.Fatalf("got %v want %v", v, archetype.Identifier(33))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
