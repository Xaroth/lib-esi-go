package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/xaroth/lib-esi-go/internal/generate/cmdutil"
	"github.com/xaroth/lib-esi-go/internal/generate/commonmodels"
	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

func main() {
	os.Exit(run())
}

func run() int {
	specFlags := cmdutil.RegisterSpecFlags(flag.CommandLine, openapi.DefaultSpecURL, "common")
	flag.Parse()

	compatDate, _ := cmdutil.CompatDateFromArgs(flag.Args(), false)
	if err := cmdutil.ValidateAndLoadCompatDate(compatDate); err != nil {
		fmt.Fprintf(os.Stderr, "invalid compatibility date %q: %v\n", compatDate, err)
		return 1
	}

	outDir, err := cmdutil.OutputDir(*specFlags.Out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "module root: %v\n", err)
		return 1
	}

	ctx := context.Background()
	spec, err := cmdutil.LoadSpec(ctx, compatDate, specFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load spec: %v\n", err)
		return 1
	}

	models, err := commonmodels.ModelsFromSpec(spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "models: %v\n", err)
		return 1
	}

	written, err := commonmodels.WritePackages(outDir, models, *specFlags.Check)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		return 1
	}

	if *specFlags.Check {
		fmt.Printf("check ok: %d packages in %s\n", len(models), outDir)
	} else {
		fmt.Printf("wrote %d files for %d packages to %s (compatibility-date: %s)\n", written, len(models), outDir, compatDate)
	}
	return 0
}
