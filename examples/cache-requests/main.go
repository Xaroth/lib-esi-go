package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "github.com/glebarez/go-sqlite"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/examples/shared"
	"github.com/xaroth/lib-esi-go/middleware"
	"github.com/xaroth/lib-esi-go/middleware/cache"
	"github.com/xaroth/lib-esi-go/request"
	"github.com/xaroth/lib-esi-go/transport"
)

type Output struct {
	Players       int64     `json:"players"`
	ServerVersion string    `json:"server_version"`
	StartTime     time.Time `json:"start_time"`
	VIP           bool      `json:"vip"`
}

var GetServerStatus = request.CreateStatic[*Output](http.MethodGet, "/status")

var requestMiddleware middleware.Middleware = func(next http.RoundTripper) http.RoundTripper {
	// This middleware only gets triggered if the request is not found in cache.
	// This means that the request headers and response status will not be shown if we are fetching directly from cache.
	// The
	return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
		req = req.Clone(req.Context())

		fmt.Printf("Request:\n")
		fmt.Printf("  If-Modified-Since: %s\n", req.Header.Get("If-Modified-Since"))
		fmt.Printf("  If-None-Match: %s\n", req.Header.Get("If-None-Match"))

		resp, err := next.RoundTrip(req)

		fmt.Printf("Response: %s\n", resp.Status)
		fmt.Printf("  Date: %s\n", resp.Header.Get("Date"))
		fmt.Printf("  Last-Modified: %s\n", resp.Header.Get("Last-Modified"))
		fmt.Printf("  ETag: %s\n", resp.Header.Get("ETag"))
		fmt.Printf("  Expires: %s\n", resp.Header.Get("Expires"))
		fmt.Printf("  Cache-Control: %s\n", resp.Header.Get("Cache-Control"))

		return resp, err
	})
}

func main() {
	ctx := context.Background()

	client := &http.Client{
		Transport: transport.New(
			"lib-esi-go examples", "0.0.0", []string{}, defaults.CompatibilityDate,
			transport.WithMiddleware(cache.Middleware("./cache.sqlite")),
			transport.WithMiddleware(requestMiddleware),
		),
	}

	resp, err := GetServerStatus(ctx, client)
	if err != nil {
		shared.Error("failed to get server status: %v", err)
	}

	if shared.ShowErrors(resp) {
		return
	}

	if resp.Data == nil {
		shared.Error("no data returned from server status")
	}

	fmt.Printf("Cache response:\n")
	fmt.Printf("  Status: %s\n", resp.Status)
	fmt.Printf("  Date: %s\n", resp.Header.Get("Date"))
	fmt.Printf("  Expires: %s\n", resp.Header.Get("Expires"))
	fmt.Printf("  Last-Modified: %s\n", resp.Header.Get("Last-Modified"))
	fmt.Printf("  ETag: %s\n", resp.Header.Get("ETag"))
	fmt.Printf("  Cache Status: %s\n", resp.Header.Get("X-Httpcache-Status"))
	fmt.Printf("Server status:\n")
	fmt.Printf("  Players: %d\n", resp.Data.Players)
	fmt.Printf("  Server version: %s\n", resp.Data.ServerVersion)
	fmt.Printf("  Start time: %s\n", resp.Data.StartTime)
	fmt.Printf("  VIP: %t\n", resp.Data.VIP)
}
