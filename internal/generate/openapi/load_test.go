package openapi_test

import (
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

func TestValidateCompatibilityDate(t *testing.T) {
	if err := openapi.ValidateCompatibilityDate("2026-02-01"); err != nil {
		t.Fatal(err)
	}
	if err := openapi.ValidateCompatibilityDate("bad"); err == nil {
		t.Fatal("expected error")
	}
}
