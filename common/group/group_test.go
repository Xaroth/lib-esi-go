package group_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/group"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `1559`
	var v group.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != group.Identifier(1559) {
		t.Fatalf("got %v want %v", v, group.Identifier(1559))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
