package item_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/item"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `1000000000001`
	var v item.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != item.Identifier(1000000000001) {
		t.Fatalf("got %v want %v", v, item.Identifier(1000000000001))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
