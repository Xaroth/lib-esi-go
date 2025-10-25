package defaults

import (
	"fmt"
	"net/url"
	"time"
)

const (
	CompatibilityDate = "2025-12-27"
	RequestTimeout    = 10 * time.Second
)

var (
	UserAgent = fmt.Sprintf("lib-esi-go/%s (+https://github.com/xaroth/lib-esi-go)", CompatibilityDate)
	Host, _   = url.Parse("https://esi.evetech.net")
)
