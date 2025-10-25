package main

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

var (
	compatibilityOptions = codegen.CompatibilityOptions{
		AlwaysPrefixEnumValues: true,
	}
)

func generateModels(ctx context.Context, cfg *Config, spec *openapi3.T) (string, error) {
	ccfg := codegen.Configuration{
		PackageName:   cfg.PackageName,
		Compatibility: compatibilityOptions,
		Generate: codegen.GenerateOptions{
			Models: true,
		},
		OutputOptions: codegen.OutputOptions{
			IncludeTags: cfg.Tags,
			SkipPrune:   true,
		},
	}
	out, err := codegen.Generate(spec, ccfg)
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}
	return out, nil
}

func generateClient(ctx context.Context, cfg *Config, spec *openapi3.T) (string, error) {
	ccfg := codegen.Configuration{
		PackageName:   cfg.PackageName,
		Compatibility: compatibilityOptions,
		Generate: codegen.GenerateOptions{
			Client: true,
		},
		OutputOptions: codegen.OutputOptions{
			IncludeTags: cfg.Tags,
		},
	}
	out, err := codegen.Generate(spec, ccfg)
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}
	return out, nil
}

func generateOutput(ctx context.Context, cfg *Config, spec *openapi3.T) error {

	ccfg := codegen.Configuration{
		PackageName:   cfg.PackageName,
		Compatibility: compatibilityOptions,
		Generate: codegen.GenerateOptions{
			Client: true,
			Models: true,
		},
		OutputOptions: codegen.OutputOptions{
			IncludeTags: cfg.Tags,
		},
	}

	_, err := codegen.Generate(spec, ccfg)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	return nil
}
