package corporation_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/corporation"
)

func TestIdentifier_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `98777771`
	var v corporation.Identifier
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != corporation.Identifier(98777771) {
		t.Fatalf("got %v want %v", v, corporation.Identifier(98777771))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
