package requestgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

// BuildPackage constructs a PackageModel from an operation.
func BuildPackage(op Operation, spec *openapi.Spec, cfg Config) (PackageModel, error) {
	resolver := openapi.NewResolver(spec)
	mapper := NewTypeMapper(cfg, spec)

	var (
		inputFields  []StructField
		inputNested  []StructDef
	)
	inSB := newStructBuilder(mapper, resolver)

	for _, pref := range op.Spec.Parameters {
		p, compName, err := resolver.ResolveParameterRef(pref)
		if err != nil {
			return PackageModel{}, fmt.Errorf("%s: %w", op.OperationID, err)
		}
		if compName != "" && isIgnoredParameter(compName) {
			continue
		}
		if p.Schema == nil {
			continue
		}
		goType, commonName, err := mapper.MapSchemaRef(*p.Schema, p.Required)
		if err != nil {
			return PackageModel{}, fmt.Errorf("%s parameter %q: %w", op.OperationID, p.Name, err)
		}
		tagKey := p.In
		if tagKey == "header" {
			tagKey = "header"
		}
		inputFields = append(inputFields, StructField{
			Name:   FieldNameFromWire(p.Name, commonName),
			Type:   goType,
			TagKey: tagKey,
			TagVal: p.Name,
		})
	}

	if op.Spec.RequestBody != nil {
		mt, ok := op.Spec.RequestBody.Content["application/json"]
		if !ok {
			return PackageModel{}, fmt.Errorf("%s: unsupported request body content type", op.OperationID)
		}
		bodyFields, nested, err := buildRequestBodyFields(mt.Schema, inSB)
		if err != nil {
			return PackageModel{}, fmt.Errorf("%s request body: %w", op.OperationID, err)
		}
		inputFields = append(inputFields, bodyFields...)
		inputNested = append(inputNested, nested...)
	}

	outputFields, outputNested, outputType, noOutputFile, outputAlias, outputAliasGoType, needsTime, err := buildResponse(op, resolver, mapper)
	if err != nil {
		return PackageModel{}, err
	}

	static := len(inputFields) == 0

	return PackageModel{
		PackageName:   PackageNameFromOperationID(op.OperationID),
		Method:        op.Method,
		Path:          op.Path,
		InputFields:   inputFields,
		InputNested:   inputNested,
		OutputFields:  outputFields,
		OutputNested:  outputNested,
		OutputType:    outputType,
		OutputAlias:       outputAlias,
		OutputAliasGoType: outputAliasGoType,
		NoOutputFile:      noOutputFile,
		Static:        static,
		NeedsTime:     needsTime,
	}, nil
}

func buildRequestBodyFields(schema openapi.SchemaRef, sb *structBuilder) ([]StructField, []StructDef, error) {
	s, _, err := sb.resolver.ResolveSchemaRef(schema)
	if err != nil {
		return nil, nil, err
	}

	if s.Type == "array" {
		itemRef := schema
		if schema.Items != nil {
			itemRef = *schema.Items
		} else if s.Items != nil {
			itemRef = *s.Items
		}
		itemSchema, itemName, err := sb.resolver.ResolveSchemaRef(itemRef)
		if err != nil {
			return nil, nil, err
		}
		if variants := oneOfVariants(itemSchema, itemRef); len(variants) > 0 {
			goType, _, err := sb.mapOneOf(variants, "body", itemName, true, true)
			if err != nil {
				return nil, nil, err
			}
			return []StructField{{
				Name:   BodyFieldName("", itemName),
				Type:   goType,
				TagKey: "body",
				TagVal: "json",
			}}, sb.nested, nil
		}
		if isObjectSchema(itemSchema) {
			name := objectStructName(itemName, "body")
			fields, _, err := sb.buildFields(itemSchema, itemName, itemRef, "json")
			if err != nil {
				return nil, nil, err
			}
			sb.registerStruct(name, fields)
			return []StructField{{
				Name:   BodyFieldName("", itemName),
				Type:   GoType{Type: "[]" + name},
				TagKey: "body",
				TagVal: "json",
			}}, sb.nested, nil
		}
		goType, itemCommon, err := sb.mapper.MapSchemaRef(itemRef, true)
		if err != nil {
			return nil, nil, err
		}
		goType.Type = "[]" + strings.TrimPrefix(goType.Type, "*")
		fieldName := BodyFieldName("", itemCommon)
		return []StructField{{
			Name:   fieldName,
			Type:   goType,
			TagKey: "body",
			TagVal: "json",
		}}, nil, nil
	}

	if s.Type == "object" {
		wires := make([]string, 0, len(s.Properties))
		for wire := range s.Properties {
			wires = append(wires, wire)
		}
		sort.Strings(wires)

		var fields []StructField
		for _, wire := range wires {
			propRef := s.Properties[wire]
			required := containsString(s.Required, wire)
			goType, commonName, _, err := sb.mapField(propRef, wire, required)
			if err != nil {
				return nil, nil, fmt.Errorf("property %q: %w", wire, err)
			}
			fields = append(fields, StructField{
				Name:   FieldNameFromWire(wire, commonName),
				Type:   goType,
				TagKey: "body",
				TagVal: "json",
			})
		}
		return fields, sb.nested, nil
	}

	goType, commonName, err := sb.mapper.MapSchemaRef(schema, true)
	if err != nil {
		return nil, nil, err
	}
	return []StructField{{
		Name:   FieldNameFromWire("body", commonName),
		Type:   goType,
		TagKey: "body",
		TagVal: "json",
	}}, nil, nil
}

func schemaFromSchema(s openapi.Schema) openapi.SchemaRef {
	return openapi.SchemaRef{
		Ref:        s.Ref,
		Type:       s.Type,
		Format:     s.Format,
		Enum:       s.Enum,
		Items:      s.Items,
		Properties: s.Properties,
		Required:   s.Required,
		OneOf:      s.OneOf,
	}
}

func buildResponse(op Operation, resolver *openapi.Resolver, mapper *TypeMapper) ([]StructField, []StructDef, string, bool, string, GoType, bool, error) {
	for _, code := range successResponseCodes(op.Method) {
		resp, ok := op.Spec.Responses[code]
		if !ok {
			continue
		}
		mt, ok := resp.Content["application/json"]
		if !ok {
			continue
		}
		return buildJSONResponse(op, mt.Schema, resolver, mapper)
	}

	if _, ok := op.Spec.Responses["204"]; ok {
		return nil, nil, "struct{}", true, "", GoType{}, false, nil
	}

	return nil, nil, "", false, "", GoType{}, false, fmt.Errorf("%s: no success response with application/json or 204", op.OperationID)
}

func successResponseCodes(method string) []string {
	if method == "POST" {
		return []string{"201", "200"}
	}
	return []string{"200"}
}

func buildJSONResponse(op Operation, schemaRef openapi.SchemaRef, resolver *openapi.Resolver, mapper *TypeMapper) ([]StructField, []StructDef, string, bool, string, GoType, bool, error) {
	schema, schemaName, err := resolver.ResolveSchemaRef(schemaRef)
	if err != nil {
		return nil, nil, "", false, "", GoType{}, false, err
	}

	sb := newStructBuilder(mapper, resolver)

	if schema.Type == "array" {
		if schema.Items == nil {
			return nil, nil, "", false, "", GoType{}, false, fmt.Errorf("%s: array response without items", op.OperationID)
		}
		itemSchema, itemName, err := resolver.ResolveSchemaRef(*schema.Items)
		if err != nil {
			return nil, nil, "", false, "", GoType{}, false, err
		}
		if schemaYieldsStruct(itemSchema, itemName, *schema.Items) {
			fields, needsTime, err := sb.buildFields(itemSchema, itemName, *schema.Items, "json")
			if err != nil {
				return nil, nil, "", false, "", GoType{}, false, err
			}
			return fields, sb.nested, "[]*Output", false, "", GoType{}, needsTime, nil
		}
		goType, _, err := mapper.MapSchemaRef(*schema.Items, true)
		if err != nil {
			return nil, nil, "", false, "", GoType{}, false, err
		}
		elem := strings.TrimPrefix(goType.Type, "*")
		alias := "[]" + elem
		goType.Type = alias
		return nil, nil, "Output", false, alias, goType, strings.Contains(elem, "time.Time"), nil
	}

	if !schemaYieldsStruct(schema, schemaName, schemaRef) {
		goType, _, err := mapper.MapSchemaRef(schemaRef, true)
		if err != nil {
			return nil, nil, "", false, "", GoType{}, false, err
		}
		alias := strings.TrimPrefix(goType.Type, "*")
		goType.Type = alias
		return nil, nil, "Output", false, alias, goType, strings.Contains(alias, "time.Time"), nil
	}

	fields, needsTime, err := sb.buildFields(schema, schemaName, schemaRef, "json")
	if err != nil {
		return nil, nil, "", false, "", GoType{}, false, err
	}
	return fields, sb.nested, "*Output", false, "", GoType{}, needsTime, nil
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
