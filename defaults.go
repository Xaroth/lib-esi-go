package defaults

import (
	"fmt"
	"net/url"
	"time"
)

const (
	Language          = "en"
	Tenant            = "tranquility"
	CompatibilityDate = "2026-07-17"
	RequestTimeout    = 10 * time.Second

	Tier = "live"
)

var UserAgent = fmt.Sprintf("lib-esi-go/%s (+https://github.com/xaroth/lib-esi-go)", CompatibilityDate)

var (
	TieredHosts = map[string]*url.URL{
		"live": asUrl("https://esi.evetech.net"),
		"test": asUrl("https://esi-test.evetech.net"),
		"dev":  asUrl("https://esi-dev.evetech.net"),
	}
	Host = TieredHosts["live"]
)

func asUrl(u string) *url.URL {
	url, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return url
}
