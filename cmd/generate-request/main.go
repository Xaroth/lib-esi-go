package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/xaroth/lib-esi-go/internal/generate/cmdutil"
	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
	"github.com/xaroth/lib-esi-go/internal/generate/requestgen"
)

func main() {
	os.Exit(run())
}

func run() int {
	specFlags := cmdutil.RegisterSpecFlags(flag.CommandLine, openapi.DefaultSpecURL, ".")
	flagLib := flag.String("lib", "github.com/xaroth/lib-esi-go", "module path for request and common model imports")
	flagCommon := flag.String("common", "common", "path suffix after -lib for common model imports")
	flag.Parse()

	compatDate, selectors := cmdutil.CompatDateFromArgs(flag.Args(), true)
	if err := cmdutil.ValidateAndLoadCompatDate(compatDate); err != nil {
		fmt.Fprintf(os.Stderr, "invalid compatibility date %q: %v\n", compatDate, err)
		return 1
	}
	if len(selectors) == 0 {
		fmt.Fprintln(os.Stderr, "usage: generate-request [compatibility-date|LIBRARY] <operation|ALL_PATHS>...")
		return 1
	}

	outDir, err := cmdutil.OutputDir(*specFlags.Out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "module root: %v\n", err)
		return 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), cmdutil.RequestSpecTimeout)
	defer cancel()

	spec, err := cmdutil.LoadSpec(ctx, compatDate, specFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load spec: %v\n", err)
		return 1
	}

	ops, err := requestgen.FindOperations(spec, selectors)
	if err != nil {
		fmt.Fprintf(os.Stderr, "find operations: %v\n", err)
		return 1
	}

	cfg := requestgen.Config{
		LibModule:    *flagLib,
		CommonSuffix: *flagCommon,
	}

	written, err := requestgen.BuildAndWrite(spec, ops, outDir, cfg, *specFlags.Check)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		return 1
	}

	if *specFlags.Check {
		fmt.Printf("check ok: %d operations in %s\n", len(ops), outDir)
	} else {
		fmt.Printf("wrote %d files for %d operations to %s (compatibility-date: %s)\n", written, len(ops), outDir, compatDate)
	}
	return 0
}
