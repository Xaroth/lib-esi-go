package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// LoadSpec loads the OpenAPI spec from inputPath or by fetching url with X-Compatibility-Date.
func LoadSpec(ctx context.Context, compatibilityDate, url, inputPath string) (*Spec, error) {
	if inputPath != "" {
		return loadSpecFile(inputPath)
	}
	if url == "" {
		url = DefaultSpecURL
	}
	return fetchSpec(ctx, url, compatibilityDate)
}

func loadSpecFile(path string) (*Spec, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var spec Spec
	if err := json.NewDecoder(f).Decode(&spec); err != nil {
		return nil, fmt.Errorf("decode spec: %w", err)
	}
	return &spec, nil
}

func fetchSpec(ctx context.Context, url, compatibilityDate string) (*Spec, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Compatibility-Date", compatibilityDate)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch spec: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("fetch spec: status %d: %s", resp.StatusCode, body)
	}

	var spec Spec
	if err := json.NewDecoder(resp.Body).Decode(&spec); err != nil {
		return nil, fmt.Errorf("decode spec: %w", err)
	}
	return &spec, nil
}

// ValidateCompatibilityDate checks date is YYYY-MM-DD.
func ValidateCompatibilityDate(date string) error {
	_, err := time.Parse(time.DateOnly, date)
	return err
}
