package writefile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-cmp/cmp"
)

// Write creates or updates path with content, or verifies content when check is true.
func Write(path string, content []byte, check bool) error {
	if check {
		existing, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("check: missing %s (would create)", path)
			}
			return err
		}
		if diff := cmp.Diff(existing, content); diff != "" {
			return fmt.Errorf("check: %s differs from generated:\n%s", path, diff)
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}
