package alliance_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/alliance"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `99000001`
	var v alliance.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != alliance.Identifier(99000001) {
		t.Fatalf("got %v want %v", v, alliance.Identifier(99000001))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
