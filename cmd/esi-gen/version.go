package main

import (
	"errors"
	"fmt"
	"runtime/debug"
)

const (
	PackageName = "esi-gen"

	LibraryName    = "lib-esi-go"
	LibraryContact = "github.com/xaroth"
)

var (
	errNoBuildInfoAvailable = errors.New("no build info available")
)

// TODO: move this to the actual library
func BuildUserAgent(vi *VersionInfo) string {
	return fmt.Sprintf("%s %s/%s (%s)", PackageName, LibraryName, vi.version, LibraryContact)
}

type VersionInfo struct {
	name      string
	version   string
	goVersion string
}

func NewVersionInfo() (*VersionInfo, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, errNoBuildInfoAvailable
	}

	return &VersionInfo{
		name:      LibraryName,
		version:   bi.Main.Version,
		goVersion: bi.GoVersion,
	}, nil
}

func (vi *VersionInfo) String() string {
	return fmt.Sprintf("%s version %s (%s)\n", PackageName, vi.version, vi.goVersion)
}
