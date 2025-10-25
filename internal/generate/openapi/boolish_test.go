package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

func TestBoolish_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{name: "json true", input: `true`, want: true},
		{name: "json false", input: `false`, want: false},
		{name: "string true", input: `"true"`, want: true},
		{name: "string false", input: `"false"`, want: false},
		{name: "invalid string", input: `"yes"`, wantErr: true},
		{name: "invalid token", input: `1`, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b openapi.Boolish
			err := json.Unmarshal([]byte(tc.input), &b)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if b.IsTrue() != tc.want {
				t.Errorf("IsTrue() = %v, want %v", b.IsTrue(), tc.want)
			}
		})
	}
}
