package requestgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/xaroth/lib-esi-go/internal/generate/commonmodels"
	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

func oneOfVariants(schema openapi.Schema, ref openapi.SchemaRef) []openapi.SchemaRef {
	if len(ref.OneOf) > 0 {
		return ref.OneOf
	}
	return schema.OneOf
}

func isOneOfSchema(schema openapi.Schema, ref openapi.SchemaRef) bool {
	return len(oneOfVariants(schema, ref)) > 0
}

func oneOfStructName(schemaName, wire string) string {
	if schemaName != "" {
		return commonmodels.TypeName(schemaName)
	}
	return singularStructName(wire)
}

func (sb *structBuilder) mapOneOf(variants []openapi.SchemaRef, wire, schemaName string, required, slice bool) (GoType, bool, error) {
	name := oneOfStructName(schemaName, wire)
	fields, needsTime, err := sb.buildOneOfFields(variants)
	if err != nil {
		return GoType{}, false, err
	}
	sb.registerStruct(name, fields)

	typ := name
	if slice {
		typ = "[]" + name
	} else if !required {
		typ = "*" + typ
	}
	return GoType{Type: typ}, needsTime, nil
}

func (sb *structBuilder) buildOneOfFields(variants []openapi.SchemaRef) ([]StructField, bool, error) {
	var fields []StructField
	needsTime := false
	for _, variantRef := range variants {
		variant, _, err := sb.resolver.ResolveSchemaRef(variantRef)
		if err != nil {
			return nil, false, err
		}
		if len(variant.Properties) != 1 {
			return nil, false, fmt.Errorf("oneOf variant: expected 1 property, got %d", len(variant.Properties))
		}
		propWires := make([]string, 0, len(variant.Properties))
		for propWire := range variant.Properties {
			propWires = append(propWires, propWire)
		}
		sort.Strings(propWires)
		for _, propWire := range propWires {
			propRef := variant.Properties[propWire]
			goType, commonName, nt, err := sb.mapField(propRef, propWire, false)
			if err != nil {
				return nil, false, fmt.Errorf("oneOf property %q: %w", propWire, err)
			}
			goType.Type = ensurePointer(goType.Type)
			fields = append(fields, StructField{
				Name:         FieldNameFromWire(propWire, commonName),
				Type:         goType,
				TagKey:       "json",
				TagVal:       propWire,
				TagOmitEmpty: true,
			})
			if nt {
				needsTime = true
			}
		}
	}
	return fields, needsTime, nil
}

func ensurePointer(typ string) string {
	if strings.HasPrefix(typ, "*") || strings.HasPrefix(typ, "[]") {
		return typ
	}
	return "*" + typ
}
