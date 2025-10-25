package commonmodels_test

import (
	"testing"

	"github.com/xaroth/lib-esi-go/internal/generate/commonmodels"
)

func TestPackageName(t *testing.T) {
	tests := []struct {
		schema string
		want   string
	}{
		{"AllianceID", "alliance"},
		{"CharacterID", "character"},
		{"TypeID", "typeid"},
		{"UUID", "uuid"},
		{"CompatibilityDate", "compatibilitydate"},
		{"SolarSystemID", "solarsystem"},
		{"ShipTreeGroupID", "shiptreegroup"},
	}
	for _, tt := range tests {
		if got := commonmodels.PackageName(tt.schema); got != tt.want {
			t.Errorf("PackageName(%q) = %q, want %q", tt.schema, got, tt.want)
		}
	}
}

func TestTypeName(t *testing.T) {
	tests := []struct {
		schema string
		want   string
	}{
		{"AllianceID", "Identifier"},
		{"UUID", "UUID"},
		{"CompatibilityDate", "CompatibilityDate"},
	}
	for _, tt := range tests {
		if got := commonmodels.TypeName(tt.schema); got != tt.want {
			t.Errorf("TypeName(%q) = %q, want %q", tt.schema, got, tt.want)
		}
	}
}

func TestEnumFieldName(t *testing.T) {
	if got := commonmodels.EnumFieldName("male", "Male"); got != "Male" {
		t.Errorf("enumFieldName with description = %q, want Male", got)
	}
	if got := commonmodels.EnumFieldName("male", ""); got != "Male" {
		t.Errorf("enumFieldName without description = %q, want Male", got)
	}
}
