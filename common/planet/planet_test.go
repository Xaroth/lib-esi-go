package planet_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/planet"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `40000002`
	var v planet.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != planet.Identifier(40000002) {
		t.Fatalf("got %v want %v", v, planet.Identifier(40000002))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
