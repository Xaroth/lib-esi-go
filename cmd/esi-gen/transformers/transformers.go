package transformers

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Spec = openapi3.T
type TransformerFn func(spec *Spec) error

var (
	transformers = []TransformerFn{}
)

func Transform(spec *openapi3.T) error {
	for _, transformer := range transformers {
		err := transformer(spec)
		if err != nil {
			return err
		}
	}
	return nil
}

func addTransformer(transformer TransformerFn) bool {
	transformers = append(transformers, transformer)
	return true
}
