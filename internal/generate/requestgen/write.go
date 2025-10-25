package requestgen

import (
	"path/filepath"

	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
	"github.com/xaroth/lib-esi-go/internal/generate/writefile"
)

// WritePackages writes generated operation packages under outDir.
func WritePackages(outDir string, packages []PackageModel, cfg Config, check bool) (written int, err error) {
	for _, pkg := range packages {
		files, err := GeneratePackage(pkg, cfg)
		if err != nil {
			return written, err
		}

		pkgDir := filepath.Join(outDir, pkg.PackageName)

		if files.Input != nil {
			if err := writefile.Write(filepath.Join(pkgDir, "input.go"), files.Input, check); err != nil {
				return written, err
			}
			written++
		}

		if files.Output != nil {
			if err := writefile.Write(filepath.Join(pkgDir, "output.go"), files.Output, check); err != nil {
				return written, err
			}
			written++
		}

		if err := writefile.Write(filepath.Join(pkgDir, "request.go"), files.Request, check); err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}

// BuildAndWrite resolves operations, builds packages, and writes files.
func BuildAndWrite(spec *openapi.Spec, ops []Operation, outDir string, cfg Config, check bool) (int, error) {
	packages := make([]PackageModel, 0, len(ops))
	for _, op := range ops {
		pkg, err := BuildPackage(op, spec, cfg)
		if err != nil {
			return 0, err
		}
		packages = append(packages, pkg)
	}
	return WritePackages(outDir, packages, cfg, check)
}
