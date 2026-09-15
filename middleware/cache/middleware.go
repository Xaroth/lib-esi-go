package cache

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/xaroth/lib-esi-go/middleware"

	"github.com/bartventer/httpcache"
)

// Middleware caches responses in a sqlite database at path. Path may carry driver
// options as a query string, e.g. ":memory:/?table=custom_cache".
func Middleware(path string) middleware.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return httpcache.NewTransport(
			DSN(path),
			httpcache.WithUpstream(next),
		)
	}
}

// DSN builds the httpcache DSN for a sqlite path. The path travels as a query
// parameter so absolute paths with drive letters and backslashes survive URL parsing.
func DSN(path string) string {
	p, rawQuery, _ := strings.Cut(path, "?")
	p = strings.TrimSuffix(p, "/")

	query, err := url.ParseQuery(rawQuery)
	if err != nil {
		query = url.Values{}
	}
	if p != "" {
		query.Set("path", p)
	}
	return CacheDriverName + "://?" + query.Encode()
}
