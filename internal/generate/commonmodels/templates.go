package commonmodels

import (
	"bytes"
	"embed"
	"fmt"
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

// ParsedTemplates returns the embedded code generation templates.
func ParsedTemplates() (*template.Template, error) {
	templatesOnce.Do(func() {
		templates, templatesErr = template.ParseFS(templateFS, "templates/*.tmpl")
	})
	return templates, templatesErr
}

// ExecuteTemplate renders a named template with data.
func ExecuteTemplate(name string, data any) (string, error) {
	tmpl, err := ParsedTemplates()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("template %q: %w", name, err)
	}
	return buf.String(), nil
}

type enumFieldData struct {
	Name      string
	WireValue string
}

type enumTemplateData struct {
	Package    string
	TypeName   string
	VarName    string
	ModulePath string
	Fields     []enumFieldData
}

type testCaseData struct {
	TypeName       string
	Index          int
	InputJSON      string
	QualifiedType  string
	WantConv       string
}

type testTemplateData struct {
	TestPackage string
	ImportPath  string
	Cases       []testCaseData
}
