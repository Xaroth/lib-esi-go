package requestgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/xaroth/lib-esi-go/internal/generate/commonmodels"
	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

// Config controls import paths for generated code.
type Config struct {
	LibModule    string // e.g. github.com/xaroth/lib-esi-go
	CommonSuffix string // e.g. common
}

func (c Config) requestImport() string {
	return c.LibModule + "/request"
}

func (c Config) commonImport(schemaName string) string {
	pkg := commonmodels.PackageName(schemaName)
	return c.LibModule + "/" + c.CommonSuffix + "/" + pkg
}

// GoType describes a mapped Go type and optional import.
type GoType struct {
	Type    string // e.g. alliance.Identifier, *int64
	Import  string // full import path if needed
	Package string // short package name for qualified type
}

// TypeMapper maps OpenAPI schemas to Go types.
type TypeMapper struct {
	cfg           Config
	resolver      *openapi.Resolver
	commonSchemas map[string]bool
}

func NewTypeMapper(cfg Config, spec *openapi.Spec) *TypeMapper {
	names := openapi.CommonSchemaNames(spec)
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	return &TypeMapper{
		cfg:           cfg,
		resolver:      openapi.NewResolver(spec),
		commonSchemas: set,
	}
}

func (m *TypeMapper) MapSchemaRef(ref openapi.SchemaRef, required bool) (GoType, string, error) {
	schema, schemaName, err := m.resolver.ResolveSchemaRef(ref)
	if err != nil {
		return GoType{}, "", err
	}
	return m.MapSchema(schema, schemaName, required)
}

func (m *TypeMapper) MapSchema(schema openapi.Schema, schemaName string, required bool) (GoType, string, error) {
	if schemaName != "" && m.commonSchemas[schemaName] {
		pkg := commonmodels.PackageName(schemaName)
		typeName := commonmodels.TypeName(schemaName)
		imp := m.cfg.commonImport(schemaName)
		qual := pkg + "." + typeName
		if len(schema.Enum) > 0 {
			qual = pkg + ".Type"
		}
		if !required {
			qual = "*" + qual
		}
		return GoType{Type: qual, Import: imp, Package: pkg}, schemaName, nil
	}

	if schema.Type == "" && schemaName == "" {
		return GoType{Type: "json.RawMessage"}, "", nil
	}

	if schema.Type == "array" {
		if schema.Items == nil {
			return GoType{}, "", fmt.Errorf("array without items")
		}
		itemSchema, itemName, err := m.resolver.ResolveSchemaRef(*schema.Items)
		if err != nil {
			return GoType{}, "", err
		}
		elem, _, err := m.MapSchema(itemSchema, itemName, true)
		if err != nil {
			return GoType{}, "", err
		}
		typ := "[]" + strings.TrimPrefix(elem.Type, "*")
		return GoType{Type: typ, Import: elem.Import, Package: elem.Package}, "", nil
	}

	base, err := m.primitiveType(schema)
	if err != nil {
		return GoType{}, "", err
	}
	if !required {
		base = "*" + base
	}
	return GoType{Type: base}, "", nil
}

func (m *TypeMapper) primitiveType(schema openapi.Schema) (string, error) {
	switch schema.Type {
	case "integer":
		switch schema.Format {
		case "int32":
			return "int32", nil
		default:
			return "int64", nil
		}
	case "number":
		return "float64", nil
	case "boolean":
		return "bool", nil
	case "string":
		if schema.Format == "date-time" {
			return "time.Time", nil
		}
		return "string", nil
	default:
		return "", fmt.Errorf("unsupported schema type %q format %q", schema.Type, schema.Format)
	}
}

// StructField is a generated struct field.
type StructField struct {
	Name         string
	Type         GoType
	TagKey       string // path, query, header, json
	TagVal       string
	TagOmitEmpty bool // append ,omitempty to JSON tag (oneOf unions)
}

// PackageModel is everything needed to render one operation package.
type PackageModel struct {
	PackageName   string
	Method        string
	Path          string
	InputFields   []StructField
	InputNested   []StructDef
	OutputFields  []StructField
	OutputNested  []StructDef
	OutputType    string // Output, *Output, []*Output, or struct{} for request generics
	OutputAlias   string // when set, output.go is `type Output = <OutputAlias>`
	OutputAliasGoType GoType // imports for OutputAlias
	NoOutputFile  bool   // true only for 204 No Content (no output.go)
	Static        bool
	NeedsTime     bool
	RequiredScopes []string
}

func collectImports(fields []StructField, cfg Config, needsTime bool) []string {
	seen := map[string]bool{
		"net/http":              true,
		cfg.requestImport():     true,
	}
	if needsTime {
		seen["time"] = true
	}
	for _, f := range fields {
		if f.Type.Import != "" {
			seen[f.Type.Import] = true
		}
		if strings.Contains(f.Type.Type, "json.RawMessage") {
			seen["encoding/json"] = true
		}
	}
	imports := make([]string, 0, len(seen))
	for imp := range seen {
		imports = append(imports, imp)
	}
	sort.Strings(imports)
	return imports
}
