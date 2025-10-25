package bloodline_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/bloodline"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `1`
	var v bloodline.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != bloodline.Identifier(1) {
		t.Fatalf("got %v want %v", v, bloodline.Identifier(1))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
