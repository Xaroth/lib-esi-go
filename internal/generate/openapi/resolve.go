package openapi

import (
	"fmt"
	"strings"
)

// Resolver resolves JSON Schema and parameter $ref pointers within a spec.
type Resolver struct {
	spec *Spec
}

func NewResolver(spec *Spec) *Resolver {
	return &Resolver{spec: spec}
}

// ResolveParameterRef resolves a parameter (inline or #/components/parameters/Name).
func (r *Resolver) ResolveParameterRef(ref ParameterRef) (Parameter, string, error) {
	if ref.Ref != "" {
		name, err := componentName(ref.Ref, "parameters")
		if err != nil {
			return Parameter{}, "", err
		}
		p, ok := r.spec.Components.Parameters[name]
		if !ok {
			return Parameter{}, "", fmt.Errorf("parameter %q not found", name)
		}
		return p, name, nil
	}
	return Parameter{
		Name:        ref.Name,
		In:          ref.In,
		Required:    ref.Required,
		Description: ref.Description,
		Schema:      ref.Schema,
	}, "", nil
}

// ResolveSchemaRef resolves a schema ref to a Schema and component name (empty if inline).
func (r *Resolver) ResolveSchemaRef(ref SchemaRef) (Schema, string, error) {
	if ref.Ref == "" {
		return schemaFromRef(ref), "", nil
	}
	name, err := componentName(ref.Ref, "schemas")
	if err != nil {
		return Schema{}, "", err
	}
	s, ok := r.spec.Components.Schemas[name]
	if !ok {
		return Schema{}, "", fmt.Errorf("schema %q not found", name)
	}
	return s, name, nil
}

// ResolveSchemaRefDeep follows schema $ref chains and merges inline fields on the last ref.
func (r *Resolver) ResolveSchemaRefDeep(ref SchemaRef) (Schema, string, error) {
	schema, name, err := r.ResolveSchemaRef(ref)
	if err != nil {
		return Schema{}, "", err
	}
	// Parameter schemas may nest $ref only in Ref field of SchemaRef wrapper.
	if ref.Ref != "" && ref.SchemaRefNested() {
		merged := schemaFromRef(ref)
		if merged.Type != "" {
			schema.Type = merged.Type
		}
		if merged.Format != "" {
			schema.Format = merged.Format
		}
	}
	return schema, name, nil
}

func (ref SchemaRef) SchemaRefNested() bool {
	return ref.Type != "" || ref.Format != "" || len(ref.Enum) > 0
}

func schemaFromRef(ref SchemaRef) Schema {
	var items *SchemaRef
	if ref.Items != nil {
		items = ref.Items
	}
	return Schema{
		Ref:               ref.Ref,
		Type:              ref.Type,
		Format:            ref.Format,
		Enum:              ref.Enum,
		Items:             items,
		Properties:        ref.Properties,
		Required:          ref.Required,
		OneOf:             ref.OneOf,
		Examples:          ref.Examples,
		XCommonModel:      ref.XCommonModel,
		XEnumDescriptions: ref.XEnumDescriptions,
	}
}

func componentName(ref, kind string) (string, error) {
	prefix := "#/components/" + kind + "/"
	if !strings.HasPrefix(ref, prefix) {
		return "", fmt.Errorf("unsupported ref %q", ref)
	}
	return strings.TrimPrefix(ref, prefix), nil
}

// CommonSchemaNames returns schema names marked x-common-model.
func CommonSchemaNames(spec *Spec) []string {
	if spec == nil {
		return nil
	}
	names := make([]string, 0)
	for name, schema := range spec.Components.Schemas {
		if schema.XCommonModel.IsTrue() {
			names = append(names, name)
		}
	}
	return names
}
