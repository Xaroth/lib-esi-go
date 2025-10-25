package dungeon_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/dungeon"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `12367`
	var v dungeon.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != dungeon.Identifier(12367) {
		t.Fatalf("got %v want %v", v, dungeon.Identifier(12367))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
