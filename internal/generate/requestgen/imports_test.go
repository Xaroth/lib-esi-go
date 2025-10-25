package requestgen

import (
	"strings"
	"testing"
)

func TestGroupImportBlockInSource(t *testing.T) {
	cfg := Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	src := []byte(`package p

import (
	"time"

	"github.com/xaroth/lib-esi-go/common/alliance"
	"github.com/xaroth/lib-esi-go/common/character"
)

type T struct{}
`)
	out := string(groupImportBlockInSource(src, cfg))
	if !strings.Contains(out, "common/alliance") || !strings.Contains(out, "time") {
		t.Fatalf("missing imports: %s", out)
	}
	timeIdx := strings.Index(out, "time")
	allianceIdx := strings.Index(out, "common/alliance")
	if timeIdx < 0 || allianceIdx < 0 || timeIdx > allianceIdx {
		t.Errorf("expected non-library imports before common imports:\n%s", out)
	}
	if !strings.Contains(out, "time\"\n\n\t\"github.com/xaroth/lib-esi-go/common/") {
		t.Errorf("expected blank line between import groups:\n%s", out)
	}
}

func TestFileImportsForFields_inputNoTime(t *testing.T) {
	cfg := Config{LibModule: "github.com/xaroth/lib-esi-go", CommonSuffix: "common"}
	fields := []StructField{{
		Type: GoType{Type: "alliance.Identifier", Import: "github.com/xaroth/lib-esi-go/common/alliance"},
	}}
	common, other := fileImportsForFields(fields, cfg, false)
	if len(other) != 0 {
		t.Fatalf("other=%v", other)
	}
	if len(common) != 1 {
		t.Fatalf("common=%v", common)
	}
}
