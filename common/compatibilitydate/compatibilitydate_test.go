package compatibilitydate_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/common/compatibilitydate"
)

func TestCompatibilityDate_JSONRoundTrip_0(t *testing.T) {
	inputJSON := `"2025-08-26"`
	var v compatibilitydate.CompatibilityDate
	if err := json.Unmarshal([]byte(inputJSON), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v != compatibilitydate.CompatibilityDate("2025-08-26") {
		t.Fatalf("got %v want %v", v, compatibilitydate.CompatibilityDate("2025-08-26"))
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != inputJSON {
		t.Fatalf("round-trip: got %s want %s", out, inputJSON)
	}
}
