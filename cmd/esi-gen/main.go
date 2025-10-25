package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xaroth/lib-esi-go/cmd/esi-gen/transformers"
)

func exitWithError(err error) {
	fmt.Fprint(os.Stdout, err.Error())
	os.Exit(1)
}

func writeToFile(cfg *Config, filename string, content string) error {
	outDir, err := filepath.Abs(cfg.Output)
	if err != nil {
		return fmt.Errorf("abs: %w", err)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %q: %w", outDir, err)
	}
	target := filepath.Join(outDir, filename)
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %q: %w", target, err)
	}
	return nil
}

func main() {
	ctx := context.Background()

	cfg, err := getConfig()
	if err != nil {
		exitWithError(err)
	}

	if cfg.Verbose {
		fmt.Fprint(os.Stdout, cfg.String())
	}

	spec, err := fetchOpenAPISpec(ctx, cfg)
	if err != nil {
		exitWithError(err)
	}

	if err := transformers.Transform(spec); err != nil {
		exitWithError(err)
	}

	models, err := generateModels(ctx, cfg, spec)
	if err != nil {
		exitWithError(err)
	}
	if err := writeToFile(cfg, "models.go", models); err != nil {
		exitWithError(err)
	}

	os.Exit(0)
}
