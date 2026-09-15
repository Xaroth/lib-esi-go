package cache_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/glebarez/go-sqlite"

	"github.com/xaroth/lib-esi-go/middleware/cache"
)

func TestDSN(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		path string
		want string
	}{
		{name: "memory", path: ":memory:", want: "lib-esi-go://?path=%3Amemory%3A"},
		{name: "memory with options", path: ":memory:/?table=custom_cache", want: "lib-esi-go://?path=%3Amemory%3A&table=custom_cache"},
		{name: "relative", path: "./cache.sqlite", want: "lib-esi-go://?path=.%2Fcache.sqlite"},
		{name: "unix absolute", path: "/var/cache/esi.sqlite", want: "lib-esi-go://?path=%2Fvar%2Fcache%2Fesi.sqlite"},
		{name: "windows absolute", path: `C:\Users\me\esi.sqlite`, want: "lib-esi-go://?path=C%3A%5CUsers%5Cme%5Cesi.sqlite"},
		{name: "options only", path: "?compress=false", want: "lib-esi-go://?compress=false"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			if got := cache.DSN(testCase.path); got != testCase.want {
				t.Fatalf("DSN(%q) = %q, want %q", testCase.path, got, testCase.want)
			}
		})
	}
}

// An absolute filesystem path must open on every OS, Windows drive letters included.
func TestMiddleware_absolutePath(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=60")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	// Not t.TempDir(): the transport has no Close, and Windows refuses to delete an
	// open sqlite file, which would fail the test during cleanup.
	dir, err := os.MkdirTemp("", "lib-esi-go-cache")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	path := filepath.Join(dir, "cache.sqlite")
	transport := cache.Middleware(path)(http.DefaultTransport)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	for i, want := range []string{"MISS", "HIT"} {
		resp, err := transport.RoundTrip(req)
		if err != nil {
			t.Fatalf("round trip %d: %v", i, err)
		}
		resp.Body.Close()
		if got := resp.Header.Get("X-Httpcache-Status"); got != want {
			t.Fatalf("round trip %d: X-Httpcache-Status = %q, want %q", i, got, want)
		}
	}
}
