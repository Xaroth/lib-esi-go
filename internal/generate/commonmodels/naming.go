package commonmodels

import (
	"strings"
	"unicode"
)

var packageOverrides = map[string]string{
	"type": "typeid",
}

var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true, "interface": true,
	"map": true, "package": true, "range": true, "return": true, "select": true,
	"struct": true, "switch": true, "type": true, "var": true,
}

// PackageName returns the Go package / directory name for a schema.
func PackageName(schemaName string) string {
	base := schemaName
	if endsWithLowercaseID(schemaName) {
		base = strings.TrimSuffix(schemaName, "ID")
	}
	name := goSafeLower(base)
	if override, ok := packageOverrides[name]; ok {
		return override
	}
	if goKeywords[name] {
		if override, ok := packageOverrides[name]; ok {
			return override
		}
		return name + "pkg"
	}
	return name
}

// TypeName returns the exported Go type name for a schema.
func TypeName(schemaName string) string {
	if endsWithLowercaseID(schemaName) {
		return "Identifier"
	}
	return schemaName
}

func endsWithLowercaseID(name string) bool {
	if !strings.HasSuffix(name, "ID") || len(name) < 3 {
		return false
	}
	r := rune(name[len(name)-3])
	return r >= 'a' && r <= 'z'
}

func goSafeLower(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			b.WriteRune(r - 'A' + 'a')
		} else if r >= 'a' && r <= 'z' {
			b.WriteRune(r)
		} else if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// EnumFieldName derives a Go identifier for an enum constant field.
func EnumFieldName(wireValue string, description string) string {
	if description != "" {
		return descriptionToIdentifier(description)
	}
	return wireToPascal(wireValue)
}

func descriptionToIdentifier(desc string) string {
	var b strings.Builder
	upper := true
	for _, r := range desc {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if upper {
				b.WriteRune(unicode.ToUpper(r))
				upper = false
			} else {
				b.WriteRune(r)
			}
		} else {
			upper = true
		}
	}
	if b.Len() == 0 {
		return wireToPascal(desc)
	}
	return b.String()
}

func wireToPascal(wire string) string {
	if wire == "" {
		return "Unknown"
	}
	parts := strings.FieldsFunc(wire, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	if len(parts) == 0 {
		return strings.ToUpper(wire[:1]) + wire[1:]
	}
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(p[1:])
		}
	}
	if b.Len() == 0 {
		return "Unknown"
	}
	return b.String()
}
