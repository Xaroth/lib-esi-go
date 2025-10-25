package cache_test

import (
	"bytes"
	"testing"

	"github.com/xaroth/lib-esi-go/middleware/cache"
)

func TestZstdRoundTrip(t *testing.T) {
	original := []byte(`{"hello":"world","nested":{"values":[1,2,3]}}`)
	encoded := cache.ZstdEncode(original)
	if bytes.Equal(encoded, original) {
		t.Fatal("expected compressed output to differ from input")
	}
	decoded, err := cache.ZstdDecode(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decoded, original) {
		t.Fatalf("got %q, want %q", decoded, original)
	}
}
