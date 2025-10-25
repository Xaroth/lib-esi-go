package transport_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/transport"
)

func TestNew_setsDefaultHeaders(t *testing.T) {
	var captured http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	rt := transport.New("TestApp", "1.2.3", []string{"https://example.com"}, defaults.CompatibilityDate)
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if captured.Get("X-Tenant") != defaults.Tenant {
		t.Errorf("X-Tenant = %q, want %q", captured.Get("X-Tenant"), defaults.Tenant)
	}
	if captured.Get("Accept-Language") != defaults.Language {
		t.Errorf("Accept-Language = %q, want %q", captured.Get("Accept-Language"), defaults.Language)
	}
	if captured.Get("X-Compatibility-Date") != defaults.CompatibilityDate {
		t.Errorf("X-Compatibility-Date = %q, want %q", captured.Get("X-Compatibility-Date"), defaults.CompatibilityDate)
	}
	ua := captured.Get("User-Agent")
	if !strings.Contains(ua, "TestApp") || !strings.Contains(ua, "1.2.3") {
		t.Errorf("User-Agent = %q", ua)
	}
}

func TestNew_usesCompatibilityDateArgument(t *testing.T) {
	var captured http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	const customDate = "2025-11-15"
	rt := transport.New("TestApp", "1.0.0", nil, customDate)
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if captured.Get("X-Compatibility-Date") != customDate {
		t.Errorf("X-Compatibility-Date = %q, want %q", captured.Get("X-Compatibility-Date"), customDate)
	}
}
