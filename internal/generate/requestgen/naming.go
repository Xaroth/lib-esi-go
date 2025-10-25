package requestgen

import (
	"strings"

	"github.com/xaroth/lib-esi-go/internal/generate/commonmodels"
)

// PackageNameFromOperationID returns the Go package name (lowercase operationId).
func PackageNameFromOperationID(operationID string) string {
	var b strings.Builder
	for _, r := range operationID {
		if r >= 'A' && r <= 'Z' {
			b.WriteRune(r - 'A' + 'a')
		} else if r >= 'a' && r <= 'z' {
			b.WriteRune(r)
		} else if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	name := b.String()
	if name == "type" {
		return "typeid"
	}
	return name
}

// FieldNameFromWire converts a wire name to a Go field name.
func FieldNameFromWire(wire string, commonSchema string) string {
	if commonSchema != "" && commonmodels.TypeName(commonSchema) == "Identifier" {
		if strings.HasSuffix(wire, "_id") {
			wire = strings.TrimSuffix(wire, "_id")
		}
	}
	return snakeToPascal(wire)
}

func snakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(strings.ToLower(p[1:]))
		}
	}
	if b.Len() == 0 {
		return "Unknown"
	}
	return b.String()
}

// singularStructName derives a struct name from a plural wire name (e.g. identities → Identity).
func singularStructName(wire string) string {
	if strings.HasSuffix(wire, "ies") && len(wire) > 3 {
		wire = wire[:len(wire)-3] + "y"
	} else if strings.HasSuffix(wire, "s") && len(wire) > 1 && !strings.HasSuffix(wire, "ss") {
		wire = wire[:len(wire)-1]
	}
	return FieldNameFromWire(wire, "")
}

// BodyFieldName picks a field name for a root JSON array request body.
func BodyFieldName(description string, itemCommonSchema string) string {
	if itemCommonSchema != "" {
		base := strings.TrimSuffix(itemCommonSchema, "ID")
		if base != "" {
			return snakeToPascal(strings.ToLower(base[:1]) + base[1:]) + "s"
		}
	}
	desc := strings.ToLower(description)
	if strings.Contains(desc, "character") {
		return "Characters"
	}
	return "Body"
}
