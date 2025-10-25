package cache_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	// This enables the sqlite driver so we can use its in-memory store to test
	_ "github.com/glebarez/go-sqlite"

	"github.com/xaroth/lib-esi-go/middleware/cache"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		dsn  string
	}{
		{
			name: "success: defaults",
			dsn:  ":memory:",
		},
		{
			name: "success: with custom table name",
			dsn:  ":memory:/?table=custom_cache",
		},
		{
			name: "success: compression enabled",
			dsn:  ":memory:/?compress=true",
		},
		{
			name: "success: compression disabled",
			dsn:  ":memory:/?compress=false",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var requests int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				w.Header().Set("Cache-Control", "public, max-age=60")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{}"))
			}))
			defer server.Close()

			middleware := cache.Middleware(testCase.dsn)
			transport := middleware(http.DefaultTransport)

			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := transport.RoundTrip(req)
			if err != nil {
				t.Fatalf("failed to round trip request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}

			if resp.Header.Get("X-Httpcache-Status") != "MISS" {
				t.Fatalf("expected X-Httpcache-Status header to be 'MISS', got '%s'", resp.Header.Get("X-Httpcache-Status"))
			}

			if requests != 1 {
				t.Fatalf("expected 1 request, got %d", requests)
			}

			resp, err = transport.RoundTrip(req)
			if err != nil {
				t.Fatalf("failed to round trip request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}

			if requests != 1 {
				t.Fatalf("expected 1 request, got %d", requests)
			}

			if resp.Header.Get("X-Httpcache-Status") != "HIT" {
				t.Fatalf("expected X-Httpcache-Status header to be 'HIT', got '%s'", resp.Header.Get("X-Httpcache-Status"))
			}
		})
	}
}

func BenchmarkMiddleware(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=60")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	middleware := cache.Middleware("//:memory:")
	transport := middleware(http.DefaultTransport)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		b.Fatalf("failed to create request: %v", err)
	}

	for b.Loop() {
		resp, err := transport.RoundTrip(req)
		if err != nil {
			b.Fatalf("failed to round trip request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	}
}
