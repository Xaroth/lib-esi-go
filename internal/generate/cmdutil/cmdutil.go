package cmdutil

import (
	"context"
	"flag"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/internal/generate/openapi"
)

const (
	// CompatDateLibrary is a positional compatibility-date token that resolves to defaults.CompatibilityDate.
	CompatDateLibrary = "LIBRARY"
	// SelectorAllPaths is the sole operation selector that exports every OpenAPI operation.
	SelectorAllPaths = "ALL_PATHS"
)

var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// SpecFlags holds standard OpenAPI generator CLI flags.
type SpecFlags struct {
	URL   *string
	Input *string
	Out   *string
	Check *bool
}

// RegisterSpecFlags registers -url, -input, -out, and -check on fs.
func RegisterSpecFlags(fs *flag.FlagSet, defaultURL, defaultOut string) SpecFlags {
	return SpecFlags{
		URL:   fs.String("url", defaultURL, "OpenAPI spec URL"),
		Input: fs.String("input", "", "local OpenAPI JSON file (instead of -url)"),
		Out:   fs.String("out", defaultOut, "output directory relative to module root"),
		Check: fs.Bool("check", false, "exit 1 if generated files differ from disk"),
	}
}

// CompatDateFromArgs returns the compatibility date and remaining args.
// When onlyIfDatePattern is true, the first arg replaces the default if it matches YYYY-MM-DD or LIBRARY.
func CompatDateFromArgs(args []string, onlyIfDatePattern bool) (compatDate string, rest []string) {
	compatDate = defaults.CompatibilityDate
	rest = args
	if len(args) == 0 {
		return compatDate, rest
	}
	if onlyIfDatePattern {
		if strings.EqualFold(args[0], CompatDateLibrary) {
			return defaults.CompatibilityDate, args[1:]
		}
		if datePattern.MatchString(args[0]) {
			compatDate = args[0]
			rest = args[1:]
		}
		return compatDate, rest
	}
	if strings.EqualFold(args[0], CompatDateLibrary) {
		return defaults.CompatibilityDate, args[1:]
	}
	compatDate = args[0]
	return compatDate, args[1:]
}

// ValidateAndLoadCompatDate validates date and prints to stderr on failure.
func ValidateAndLoadCompatDate(date string) error {
	return openapi.ValidateCompatibilityDate(date)
}

// OutputDir resolves the output directory. "." uses the current working directory
// (e.g. the package directory when invoked via go generate); other values are
// joined relative to the module root.
func OutputDir(outFlag string) (string, error) {
	if outFlag == "." {
		return filepath.Abs(".")
	}
	root, err := openapi.ModuleRoot(".")
	if err != nil {
		return "", err
	}
	return filepath.Join(root, outFlag), nil
}

// LoadSpec loads the OpenAPI document using flags.
func LoadSpec(ctx context.Context, compatDate string, flags SpecFlags) (*openapi.Spec, error) {
	return openapi.LoadSpec(ctx, compatDate, *flags.URL, *flags.Input)
}

// RequestSpecTimeout is the default HTTP timeout for generate-request spec fetches.
const RequestSpecTimeout = 2 * time.Minute
