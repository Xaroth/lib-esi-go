package requestgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/xaroth/lib-esi-go/internal/generate/cmdutil"
	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

// Operation is a resolved OpenAPI operation.
type Operation struct {
	OperationID string
	Method      string
	Path        string
	Spec        openapi.Operation
}

var httpMethods = map[string]bool{
	"get": true, "post": true, "put": true, "patch": true, "delete": true, "head": true, "options": true,
}

// FindOperations resolves selectors against spec paths.
// When selectors is exactly ALL_PATHS (case-insensitive), every HTTP operation is returned.
func FindOperations(spec *openapi.Spec, selectors []string) ([]Operation, error) {
	if len(selectors) == 0 {
		return nil, fmt.Errorf("no operations specified")
	}
	if len(selectors) == 1 && strings.EqualFold(selectors[0], cmdutil.SelectorAllPaths) {
		return findAllOperations(spec)
	}
	var out []Operation
	for _, sel := range selectors {
		op, err := findOne(spec, sel)
		if err != nil {
			return nil, err
		}
		out = append(out, op)
	}
	return out, nil
}

func findAllOperations(spec *openapi.Spec) ([]Operation, error) {
	var out []Operation
	for path, item := range spec.Paths {
		for method, op := range item {
			if !httpMethods[strings.ToLower(method)] {
				continue
			}
			if op.OperationID == "" {
				continue
			}
			out = append(out, Operation{
				OperationID: op.OperationID,
				Method:      strings.ToUpper(method),
				Path:        path,
				Spec:        op,
			})
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no operations found in spec")
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].OperationID != out[j].OperationID {
			return out[i].OperationID < out[j].OperationID
		}
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return out[i].Method < out[j].Method
	})
	return out, nil
}

func findOne(spec *openapi.Spec, selector string) (Operation, error) {
	if method, path, ok := parseMethodPath(selector); ok {
		return findByMethodPath(spec, method, path)
	}
	return findByOperationID(spec, selector)
}

func parseMethodPath(selector string) (method, path string, ok bool) {
	parts := strings.Fields(selector)
	if len(parts) != 2 {
		return "", "", false
	}
	method = strings.ToLower(parts[0])
	if !httpMethods[method] {
		return "", "", false
	}
	return method, parts[1], true
}

func findByOperationID(spec *openapi.Spec, id string) (Operation, error) {
	for path, item := range spec.Paths {
		for method, op := range item {
			if !httpMethods[strings.ToLower(method)] {
				continue
			}
			if strings.EqualFold(op.OperationID, id) {
				return Operation{
					OperationID: op.OperationID,
					Method:      strings.ToUpper(method),
					Path:        path,
					Spec:        op,
				}, nil
			}
		}
	}
	return Operation{}, fmt.Errorf("operation %q not found", id)
}

func findByMethodPath(spec *openapi.Spec, method, path string) (Operation, error) {
	item, ok := spec.Paths[path]
	if !ok {
		return Operation{}, fmt.Errorf("path %q not found", path)
	}
	op, ok := item[strings.ToLower(method)]
	if !ok {
		return Operation{}, fmt.Errorf("%s %q not found", strings.ToUpper(method), path)
	}
	return Operation{
		OperationID: op.OperationID,
		Method:      strings.ToUpper(method),
		Path:        path,
		Spec:        op,
	}, nil
}
