package cache

import (
	"fmt"
	"net/http"
	"time"

	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/request"
)

// Middleware automatically detects ETag and Last-Modified headers and sends them
// with subsequent requests to save bandwidth and improve performance.
// This middleware is opt-in, and is not enabled by default.
func Middleware(storage Cache) middleware.Middleware {
	if storage == nil {
		storage = NewMemoryStorage(1000)
	}

	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {

			ctx := req.Context()
			info, ok := request.GetRequestInfo(ctx)
			// Only cache known GET requests.
			if !ok || info.Method != http.MethodGet {
				return next.RoundTrip(req)
			}

			requestKey := request.GetRequestKey(ctx)
			cacheKey := fmt.Sprintf("%s:%s", info.Path, requestKey)

			if entry, ok := storage.Get(cacheKey); ok {
				// Clone the request per RoundTripper contract
				req = req.Clone(req.Context())

				// Add conditional headers to the request
				if entry.ETag != "" {
					req.Header.Set("If-None-Match", entry.ETag)
				}
				if !entry.LastModified.IsZero() {
					req.Header.Set("If-Modified-Since", entry.LastModified.Format(http.TimeFormat))
				}
			}

			resp, err := next.RoundTrip(req)

			if err != nil {
				return nil, err
			}

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				etag := resp.Header.Get("ETag")
				if lastModifiedStr := resp.Header.Get("Last-Modified"); lastModifiedStr != "" {
					if lastModified, err := http.ParseTime(lastModifiedStr); err == nil {
						storage.Set(cacheKey, &CacheEntry{
							ETag:         etag,
							LastModified: lastModified,
						})
						return resp, nil
					}
				}

				// If we didn't have a last-modified header, assume last modified is now.
				storage.Set(cacheKey, &CacheEntry{
					ETag:         etag,
					LastModified: time.Now(),
				})
				return resp, nil
			}

			return resp, nil
		})
	}
}
