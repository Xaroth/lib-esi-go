package commonmodels

import (
	"path/filepath"

	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
	"github.com/xaroth/lib-esi-go/internal/generate/writefile"
)

// WritePackages writes generated files under outDir.
func WritePackages(outDir string, models []Model, check bool) (written int, err error) {
	moduleRoot, err := openapi.ModuleRoot(outDir)
	if err != nil {
		return 0, err
	}
	modulePath, err := openapi.ModulePath(outDir)
	if err != nil {
		return 0, err
	}
	importBase, err := openapi.ImportBase(moduleRoot, outDir)
	if err != nil {
		return 0, err
	}

	for _, m := range models {
		mainGo, testGo, err := GeneratePackage(m, modulePath, importBase)
		if err != nil {
			return written, err
		}

		pkgDir := filepath.Join(outDir, m.Package)
		mainPath := filepath.Join(pkgDir, m.Package+".go")
		if err := writefile.Write(mainPath, mainGo, check); err != nil {
			return written, err
		}
		written++

		if testGo != nil {
			testPath := filepath.Join(pkgDir, m.Package+"_test.go")
			if err := writefile.Write(testPath, testGo, check); err != nil {
				return written, err
			}
		}
	}
	return written, nil
}
