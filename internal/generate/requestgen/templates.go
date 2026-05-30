package requestgen

import (
	"bytes"
	"embed"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

var (
	templates     *template.Template
	templatesOnce sync.Once
	templatesErr  error
)

func parsedTemplates() (*template.Template, error) {
	templatesOnce.Do(func() {
		templates, templatesErr = template.ParseFS(templateFS, "templates/*.tmpl")
	})
	return templates, templatesErr
}

func executeTemplate(name string, data any) (string, error) {
	tmpl, err := parsedTemplates()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("template %q: %w", name, err)
	}
	return buf.String(), nil
}

type outputAliasTemplateData struct {
	PackageName   string
	AliasType     string
	CommonImports []string
	OtherImports  []string
	HasImports    bool
}

type fileTemplateData struct {
	PackageName   string
	RootName      string
	CommonImports []string
	OtherImports  []string
	HasImports    bool
	Fields        []StructField
	Nested        []StructDef
}

type requestTemplateData struct {
	PackageName           string
	RequestImport         string
	MethodConst           string
	PathLiteral           string
	OutputType            string
	Static                bool
	HasRequiredScopes     bool
	RequiredScopesLiteral string
}

func fileImportsForFields(fields []StructField, cfg Config, needsTime bool) (common, other []string) {
	imps := collectImports(fields, cfg, needsTime)
	var paths []string
	for _, imp := range imps {
		if imp == "net/http" || imp == cfg.requestImport() {
			continue
		}
		paths = append(paths, imp)
	}
	return importGroups(paths, cfg)
}

func methodConstFixed(method string) string {
	switch method {
	case "GET":
		return "http.MethodGet"
	case "POST":
		return "http.MethodPost"
	case "PUT":
		return "http.MethodPut"
	case "PATCH":
		return "http.MethodPatch"
	case "DELETE":
		return "http.MethodDelete"
	case "HEAD":
		return "http.MethodHead"
	case "OPTIONS":
		return "http.MethodOptions"
	default:
		return "http.MethodGet"
	}
}

func pathLiteral(path string) string {
	return strconv.Quote(path)
}

func scopesLiteral(scopes []string) string {
	parts := make([]string, len(scopes))
	for i, scope := range scopes {
		parts[i] = strconv.Quote(scope)
	}
	return strings.Join(parts, ", ")
}
