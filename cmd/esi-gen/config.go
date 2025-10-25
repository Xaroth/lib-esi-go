package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	defaultTier    = "live"
	defaultOutput  = ""
	defaultPackage = "esi"
)

var (
	tier    string = defaultTier
	compat  string = ""
	out     string = defaultOutput
	pkg     string = defaultPackage
	verbose bool   = false

	showUsage   bool
	showVersion bool
)

// Config captures CLI settings for esi-gen.
type Config struct {
	Tier              string
	BaseURL           string
	CompatibilityDate string
	Output            string
	PackageName       string
	Verbose           bool
	VersionInfo       *VersionInfo
	Tags              []string
}

func parseFlags() {
	flag.StringVar(&tier, "tier", defaultTier, "ESI tier (dev|test|live)")
	flag.StringVar(&compat, "compatibility-date", "", "Compatibility date (YYYY-MM-DD)")
	flag.StringVar(&out, "o", defaultOutput, "Output file or directory (like oapi-codegen -o)")
	flag.StringVar(&pkg, "package", defaultPackage, "Go package name (like oapi-codegen -package)")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")

	flag.BoolVar(&showVersion, "version", false, "Show version and exit")

	flag.BoolVar(&showUsage, "help", false, "Show help and exit")
	flag.BoolVar(&showUsage, "h", false, "Show help and exit")

	flag.Parse()
}

// getConfig builds a Config from flags and positional args (tags).
func getConfig() (*Config, error) {
	parseFlags()

	vi, err := NewVersionInfo()
	if err != nil {
		return nil, err
	}

	if showUsage {
		flag.Usage()
		os.Exit(0)
		return nil, nil
	}

	if showVersion {
		fmt.Println(vi.String())
		os.Exit(0)
		return nil, nil
	}

	if compat == "" {
		compat = time.Now().UTC().Add(time.Hour * -11).Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", compat); err != nil {
		return nil, fmt.Errorf("invalid compatibility date %q: %w", compat, err)
	}

	var baseURL string
	switch tier {
	case "live":
		baseURL = "https://esi.evetech.net"
	case "dev", "test":
		baseURL = fmt.Sprintf("https://esi-%s.evetech.net", tier)
	default:
		return nil, fmt.Errorf("invalid tier %q: must be one of dev, test, live", tier)
	}

	cfg := &Config{
		Tier:              tier,
		BaseURL:           baseURL,
		CompatibilityDate: compat,
		Output:            out,
		PackageName:       pkg,
		Verbose:           verbose,
		VersionInfo:       vi,
		Tags:              flag.Args(),
	}
	return cfg, nil
}

func (c *Config) String() string {
	sb := strings.Builder{}

	sb.WriteString("esi-gen\n")
	sb.WriteString(fmt.Sprintf("  tier: %s\n", c.Tier))
	sb.WriteString(fmt.Sprintf("  base-url: %s\n", c.BaseURL))
	sb.WriteString(fmt.Sprintf("  compatibility-date: %s\n", c.CompatibilityDate))
	sb.WriteString(fmt.Sprintf("  package: %s\n", c.PackageName))
	if c.Output != "" {
		sb.WriteString(fmt.Sprintf("  output: %s\n", c.Output))
	} else {
		sb.WriteString("  output: stdout\n")
	}
	if len(c.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("  tags: %v\n", c.Tags))
	} else {
		sb.WriteString("  tags: (all)\n")
	}
	return sb.String()
}
