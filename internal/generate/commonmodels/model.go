package commonmodels

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

// Model is a common schema ready for code generation.
type Model struct {
	SchemaName string
	Package    string
	TypeName   string
	Schema     openapi.Schema
	Examples   []any
}

// ModelsFromSpec returns all x-common-model schemas sorted by name.
func ModelsFromSpec(spec *openapi.Spec) ([]Model, error) {	if spec == nil {
		return nil, fmt.Errorf("nil spec")
	}

	names := make([]string, 0)
	for name, schema := range spec.Components.Schemas {
		if schema.XCommonModel.IsTrue() {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	models := make([]Model, 0, len(names))
	for _, name := range names {
		schema := spec.Components.Schemas[name]
		examples, err := parseExamples(schema.Examples)
		if err != nil {
			return nil, fmt.Errorf("%s examples: %w", name, err)
		}
		models = append(models, Model{
			SchemaName: name,
			Package:    PackageName(name),
			TypeName:   TypeName(name),
			Schema:     schema,
			Examples:   examples,
		})
	}
	return models, nil
}

func parseExamples(raw json.RawMessage) ([]any, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var single any
	if err := json.Unmarshal(raw, &single); err != nil {
		return nil, err
	}

	switch v := single.(type) {
	case []any:
		return v, nil
	default:
		return []any{v}, nil
	}
}
