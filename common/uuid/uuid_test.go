package uuid_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/uuid"
)

func TestUUID_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `"3868eaed-8278-4cb7-9709-7d7de9c20dc7"`
	var v uuid.UUID
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != uuid.UUID("3868eaed-8278-4cb7-9709-7d7de9c20dc7") {
		t.Fatalf("got %v want %v", v, uuid.UUID("3868eaed-8278-4cb7-9709-7d7de9c20dc7"))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
