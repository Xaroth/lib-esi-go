package requestgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/xaroth/lib-esi-go/internal/generate/commonmodels"
	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

// StructDef is a generated struct type (root or nested).
type StructDef struct {
	Name   string
	Fields []StructField
}

type structBuilder struct {
	mapper   *TypeMapper
	resolver *openapi.Resolver
	nested   []StructDef
}

func newStructBuilder(mapper *TypeMapper, resolver *openapi.Resolver) *structBuilder {
	return &structBuilder{mapper: mapper, resolver: resolver}
}

func (sb *structBuilder) buildFields(schema openapi.Schema, schemaName string, ref openapi.SchemaRef, tagKey string) ([]StructField, bool, error) {
	if schemaName != "" && schema.Type == "" {
		s, name, err := sb.resolver.ResolveSchemaRef(ref)
		if err != nil {
			return nil, false, err
		}
		schema = s
		schemaName = name
	}
	if schema.Type == "array" && schema.Items != nil {
		itemSchema, itemName, err := sb.resolver.ResolveSchemaRef(*schema.Items)
		if err != nil {
			return nil, false, err
		}
		return sb.buildFields(itemSchema, itemName, *schema.Items, tagKey)
	}

	needsTime := false
	wires := make([]string, 0, len(schema.Properties))
	for wire := range schema.Properties {
		wires = append(wires, wire)
	}
	sort.Strings(wires)

	var fields []StructField
	for _, wire := range wires {
		propRef := schema.Properties[wire]
		required := containsString(schema.Required, wire)
		goType, commonName, nt, err := sb.mapField(propRef, wire, required)
		if err != nil {
			return nil, false, fmt.Errorf("property %q: %w", wire, err)
		}
		if nt {
			needsTime = true
		}
		fields = append(fields, StructField{
			Name:        FieldNameFromWire(wire, commonName),
			Type:        goType,
			TagKey:      tagKey,
			TagVal:      wire,
			TagRequired: required && tagKey != "path",
		})
	}
	return fields, needsTime, nil
}

func (sb *structBuilder) mapField(ref openapi.SchemaRef, wire string, required bool) (GoType, string, bool, error) {
	schema, schemaName, err := sb.resolver.ResolveSchemaRef(ref)
	if err != nil {
		return GoType{}, "", false, err
	}

	if schemaName != "" && sb.mapper.commonSchemas[schemaName] {
		goType, _, err := sb.mapper.MapSchema(schema, schemaName, required)
		return goType, schemaName, strings.Contains(goType.Type, "time.Time"), err
	}

	if variants := oneOfVariants(schema, ref); len(variants) > 0 {
		goType, nt, err := sb.mapOneOf(variants, wire, schemaName, required, false)
		return goType, "", nt, err
	}

	if isObjectSchema(schema) {
		name := objectStructName(schemaName, wire)
		fields, needsTime, err := sb.buildFields(schema, schemaName, ref, "json")
		if err != nil {
			return GoType{}, "", false, err
		}
		sb.registerStruct(name, fields)
		typ := name
		if !required {
			typ = "*" + typ
		}
		return GoType{Type: typ}, "", needsTime, nil
	}

	if schema.Type == "array" && schema.Items != nil {
		itemSchema, itemName, err := sb.resolver.ResolveSchemaRef(*schema.Items)
		if err != nil {
			return GoType{}, "", false, err
		}
		if variants := oneOfVariants(itemSchema, *schema.Items); len(variants) > 0 {
			goType, nt, err := sb.mapOneOf(variants, wire, itemName, true, true)
			return goType, "", nt, err
		}
		if isObjectSchema(itemSchema) {
			name := objectStructName(itemName, wire)
			fields, needsTime, err := sb.buildFields(itemSchema, itemName, *schema.Items, "json")
			if err != nil {
				return GoType{}, "", false, err
			}
			sb.registerStruct(name, fields)
			return GoType{Type: "[]" + name}, "", needsTime, nil
		}
	}

	goType, commonName, err := sb.mapper.MapSchemaRef(ref, required)
	return goType, commonName, strings.Contains(goType.Type, "time.Time"), err
}

func (sb *structBuilder) registerStruct(name string, fields []StructField) {
	for i, s := range sb.nested {
		if s.Name == name {
			sb.nested[i].Fields = mergeStructFields(s.Fields, fields)
			return
		}
	}
	sb.nested = append(sb.nested, StructDef{Name: name, Fields: fields})
}

func mergeStructFields(existing, add []StructField) []StructField {
	if len(add) == 0 {
		return existing
	}
	seen := make(map[string]bool, len(existing))
	for _, f := range existing {
		seen[f.Name] = true
	}
	out := append([]StructField(nil), existing...)
	for _, f := range add {
		if !seen[f.Name] {
			out = append(out, f)
			seen[f.Name] = true
		}
	}
	return out
}

func isObjectSchema(schema openapi.Schema) bool {
	if schema.Type == "object" {
		return true
	}
	return len(schema.Properties) > 0
}

// schemaYieldsStruct reports whether a schema is rendered as a generated struct
// (object, oneOf union, etc.) rather than a primitive or common-model alias.
func schemaYieldsStruct(schema openapi.Schema, schemaName string, ref openapi.SchemaRef) bool {
	if isOneOfSchema(schema, ref) {
		return true
	}
	return isObjectSchema(schema) && (len(schema.Properties) > 0 || schemaName != "")
}

func objectStructName(schemaName, wire string) string {
	if schemaName != "" {
		return commonmodels.TypeName(schemaName)
	}
	return FieldNameFromWire(wire, "")
}

func allStructFields(defs []StructDef) []StructField {
	var out []StructField
	for _, d := range defs {
		out = append(out, d.Fields...)
	}
	return out
}
