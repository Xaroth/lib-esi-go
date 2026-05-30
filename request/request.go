package request

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/request/esierror"
	"github.com/xaroth/lib-esi-go/request/internal/parameters"
	"github.com/xaroth/lib-esi-go/request/internal/pattern"
)

type RequestSender interface {
	Do(req *http.Request) (*http.Response, error)
}

type Response[TOutput any] struct {
	*http.Response
	// If an error occurs, the Data field is not guaranteed to be set.
	Data TOutput

	ErrorData *esierror.ErrorData
}

type requestInfo struct {
	Method        string
	Path          string
	Pattern       pattern.Pattern
	RequiredScope []string
}

type RequestFunc[TInput any, TOutput any] func(ctx context.Context, sender RequestSender, input *TInput, opts ...RequestOption) (*Response[TOutput], error)
type StaticFunc[TOutput any] func(ctx context.Context, sender RequestSender, opts ...RequestOption) (*Response[TOutput], error)

func getSortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Create a unique string key of the request parameters.
func createRequestKey(pathParameters map[string]any, queryParameters url.Values, headerParameters http.Header) string {
	parts := make([]string, len(pathParameters)+len(queryParameters)+len(headerParameters))

	for _, key := range getSortedKeys(pathParameters) {
		parts = append(parts, key, fmt.Sprintf("%v", pathParameters[key]))
	}
	for _, key := range getSortedKeys(queryParameters) {
		values := queryParameters[key]
		parts = append(parts, key, strings.Join(values, ","))
	}
	for _, key := range getSortedKeys(headerParameters) {
		values := headerParameters[key]
		parts = append(parts, key, strings.Join(values, ","))
	}
	return strings.Join(parts, ":")
}

func Create[TInput any, TOutput any](method string, path string, opts ...CreateOption) RequestFunc[TInput, TOutput] {
	pattern, err := pattern.NewValidated[TInput](path)
	if err != nil {
		panic(err)
	}

	req := &requestInfo{
		Method:  method,
		Path:    path,
		Pattern: pattern,
	}

	for _, opt := range opts {
		opt(req)
	}

	return func(bCtx context.Context, sender RequestSender, input *TInput, opts ...RequestOption) (*Response[TOutput], error) {
		// Split the input parameters into path, query, header, and body parameters.
		pathParameters, queryParameters, headerParameters, bodyParameters, err := parameters.Extract(input)
		if err != nil {
			return nil, err
		}

		requestKey := createRequestKey(pathParameters, queryParameters, headerParameters)

		ctx := BaseContext(bCtx, req, requestKey, input)

		path, err := pattern.String(pathParameters)
		if err != nil {
			return nil, err
		}

		// Requests are always created using the default (live) API domain.
		//
		// If a different tier is is needed, you can accomplish this with the tier middleware,
		// or by using the WithTier option when creating the transport chain.
		url := defaults.Host.JoinPath(path)
		url.RawQuery = queryParameters.Encode()

		req, err := http.NewRequestWithContext(ctx, req.Method, url.String(), bodyParameters)
		if err != nil {
			return nil, err
		}

		for key, values := range headerParameters {
			req.Header.Del(key)
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		rawResp, err := sender.Do(req)
		if err != nil {
			return nil, err
		}
		defer rawResp.Body.Close()

		data, err := io.ReadAll(rawResp.Body)
		if err != nil {
			return nil, err
		}

		resp := &Response[TOutput]{
			Response: rawResp,
		}

		if resp.StatusCode >= http.StatusBadRequest {
			if errData, err := esierror.UnmarshalErrorJSON(data); err == nil {
				resp.ErrorData = errData
			}
		} else if len(data) > 0 {
			if err := json.Unmarshal(data, &resp.Data); err != nil {
				return resp, err
			}
		}

		return resp, nil
	}
}

func CreateStatic[TOutput any](method string, path string, opts ...CreateOption) StaticFunc[TOutput] {
	base := Create[struct{}, TOutput](method, path, opts...)

	return func(ctx context.Context, sender RequestSender, opts ...RequestOption) (*Response[TOutput], error) {
		return base(ctx, sender, nil, opts...)
	}
}
