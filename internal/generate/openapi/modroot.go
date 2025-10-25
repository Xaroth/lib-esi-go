package openapi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ModuleRoot finds the directory containing go.mod by walking up from dir.
func ModuleRoot(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(abs, "go.mod")); err == nil {
			return abs, nil
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return "", fmt.Errorf("go.mod not found from %s", dir)
		}
		abs = parent
	}
}

// ModulePath returns the module path declared in the go.mod nearest to dir.
func ModulePath(dir string) (string, error) {
	root, err := ModuleRoot(dir)
	if err != nil {
		return "", err
	}
	return ModulePathFromFile(filepath.Join(root, "go.mod"))
}

// ModulePathFromFile reads the module path from a go.mod file.
func ModulePathFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	for line := range strings.SplitSeq(string(data), "\n") {
		line, ok := strings.CutPrefix(strings.TrimSpace(line), "module ")
		if ok && line != "" {
			if before, _, ok := strings.Cut(line, "//"); ok {
				line = strings.TrimSpace(before)
			}
			return line, nil
		}
	}
	return "", fmt.Errorf("module directive not found in %s", path)
}

// ImportBase returns the slash-separated import path suffix from module root to outDir.
func ImportBase(moduleRoot, outDir string) (string, error) {
	rel, err := filepath.Rel(moduleRoot, outDir)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}
