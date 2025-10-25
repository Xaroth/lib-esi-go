package character_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/character"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `90000001`
	var v character.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != character.Identifier(90000001) {
		t.Fatalf("got %v want %v", v, character.Identifier(90000001))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
