package useragent

import (
	"fmt"
	"net/http"
	"strings"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/middleware"
)

// Middleware adds a User-Agent header to each request.
// This middleware is always added to the transport chain.
func Middleware(applicationName, applicationVersion string, contact ...string) middleware.Middleware {
	userAgent := buildUserAgent(applicationName, applicationVersion, contact)

	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
			req = req.Clone(req.Context())
			req.Header.Set("User-Agent", userAgent)
			return next.RoundTrip(req)
		})
	}
}

func buildUserAgent(applicationName, applicationVersion string, applicationContact []string) string {
	sb := strings.Builder{}

	if applicationName != "" {
		sb.WriteString(applicationName)
		if applicationVersion != "" {
			sb.WriteString(fmt.Sprintf("/%s", applicationVersion))
		}
		if len(applicationContact) > 0 {
			sb.WriteString(fmt.Sprintf(" (%s)", strings.Join(applicationContact, "; ")))
		}
	} else {
		sb.WriteString("Unconfigured")
	}

	sb.WriteString("; ")
	sb.WriteString(defaults.UserAgent)
	return sb.String()
}
