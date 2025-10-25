package transformers

import (
	"fmt"
	"strings"
)

func removeParameter(spec *Spec, name string) error {
	parameterRef := fmt.Sprintf("#/components/parameters/%s", name)
	for _, path := range spec.Paths.Map() {
		for _, operation := range path.Operations() {
			for i, p := range operation.Parameters {
				if p.Ref == parameterRef {
					operation.Parameters = append(operation.Parameters[:i], operation.Parameters[i+1:]...)
				}
			}
		}
	}
	delete(spec.Components.Parameters, name)
	return nil
}

func findParameter(spec *Spec, name string, in string) string {
	for paramName, param := range spec.Components.Parameters {
		if param.Value == nil {
			continue
		}

		if strings.EqualFold(param.Value.Name, name) && param.Value.In == in {
			return paramName
		}
	}
	return ""
}

var (
	removedHeaders = []string{
		// lib-esi-go is responsible for defining this header, as it should always
		// match the compatibility date used to generate the code.
		"X-Compatibility-Date",

		// Accept-Language is an odd one. Ideally we'd want to give the user the
		// authority to specify the requested language, but oapi-codegen generates
		// individual types for each operation, causing a lot of duplicate code.
		// As such, lib-esi-go will just have to deal with it in a cleaner way.
		"Accept-Language",

		"X-Tenant",
		"User-Agent",
	}
)

// Remove any input headers that lib-esi-go deals with internally.
// This ensures that these headers are not exposed in the generated code, making
// it less likely for users to accidentally set them.
var _ = addTransformer(func(spec *Spec) error {
	for _, headerName := range removedHeaders {
		fmt.Printf("Checking for header: %s\n", headerName)
		if name := findParameter(spec, headerName, "header"); name != "" {
			fmt.Printf("Removing header: %s\n", name)
			if err := removeParameter(spec, name); err != nil {
				return err
			}
		}
	}
	return nil
})
