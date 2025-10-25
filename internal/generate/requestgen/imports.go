package requestgen

import (
	"bytes"
	"sort"
	"strings"
)

// importGroups splits import paths into common model paths and everything else.
func importGroups(imports []string, cfg Config) (common, other []string) {
	prefix := cfg.LibModule + "/" + cfg.CommonSuffix + "/"
	for _, imp := range imports {
		if strings.HasPrefix(imp, prefix) {
			common = append(common, imp)
		} else {
			other = append(other, imp)
		}
	}
	sort.Strings(common)
	sort.Strings(other)
	return common, other
}

func fieldsNeedTime(fields []StructField) bool {
	for _, f := range fields {
		if strings.Contains(f.Type.Type, "time.Time") {
			return true
		}
	}
	return false
}

// groupImportBlockInSource reorders a formatted Go file so non-library imports
// (stdlib, etc.) are listed first, then a blank line, then common model imports.
func groupImportBlockInSource(src []byte, cfg Config) []byte {
	const marker = "import ("
	idx := bytes.Index(src, []byte(marker))
	if idx < 0 {
		return src
	}
	start := idx + len(marker)
	end := bytes.Index(src[start:], []byte(")"))
	if end < 0 {
		return src
	}
	end += start

	var paths []string
	for _, line := range bytes.Split(src[start:end], []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		path := string(line)
		path = strings.TrimPrefix(path, `"`)
		path = strings.TrimSuffix(path, `"`)
		if path != "" {
			paths = append(paths, path)
		}
	}

	common, other := importGroups(paths, cfg)
	if len(common) == 0 || len(other) == 0 {
		return src
	}

	var block bytes.Buffer
	block.WriteString(marker)
	block.WriteByte('\n')
	for _, imp := range other {
		block.WriteString("\t\"")
		block.WriteString(imp)
		block.WriteString("\"\n")
	}
	block.WriteByte('\n')
	for _, imp := range common {
		block.WriteString("\t\"")
		block.WriteString(imp)
		block.WriteString("\"\n")
	}
	block.WriteByte(')')

	out := make([]byte, 0, len(src)+64)
	out = append(out, src[:idx]...)
	out = append(out, block.Bytes()...)
	out = append(out, src[end+1:]...)
	return out
}
