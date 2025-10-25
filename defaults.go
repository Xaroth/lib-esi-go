package request

import (
	"fmt"
	"time"
)

const (
	CompatibilityDate = "2025-11-22"

	Host            = "https://esi.evetech.net"
	Language        = "en"
	RequestTimeout  = 10 * time.Second
	ScheduleTimeout = 900 * time.Second // 15 minutes
)

var (
	UserAgent = fmt.Sprintf("lib-esi-go/%s (+https://github.com/xaroth/lib-esi-go)", CompatibilityDate)
)
