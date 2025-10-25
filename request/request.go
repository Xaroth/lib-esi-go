package request

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/xaroth/lib-esi-go/middleware/authentication"
)

var (
	ErrInvalidValueType  = errors.New("invalid value type")
	ErrInvalidBodyType   = errors.New("invalid body type")
	ErrInvalidStatusCode = errors.New("invalid status code")

	ErrUndefinedVariable = errors.New("undefined variable")

	ErrNotModified = errors.New("not modified")
)

type RequestSender interface {
	Do(req *http.Request) (*http.Response, error)
}

type Response[TOutput any] struct {
	*http.Response
	// If an error occurs, the Data field is not guaranteed to be set.
	Data TOutput

	ErrorData *ErrorData
}

func validate[TInput any](pattern Pattern) error {
	variables := pattern.Variables()

	foundVariables := make(map[string]bool)

	var input TInput
	typ := reflect.TypeOf(input)
	for i := range typ.NumField() {
		field := typ.Field(i)
		if value, ok := field.Tag.Lookup("path"); ok {
			foundVariables[value] = true
		}
	}

	for _, variable := range variables {
		if !foundVariables[variable] {
			return fmt.Errorf("%w: %s", ErrUndefinedVariable, variable)
		}
	}
	for variable := range foundVariables {
		if !slices.Contains(variables, variable) {
			return fmt.Errorf("%w: %s", ErrExtraneousVariable, variable)
		}
	}

	return nil
}

type RequestInfo struct {
	Method        string
	Path          string
	RequiredScope string
}

type RequestFunc[TInput any, TOutput any] func(ctx context.Context, sender RequestSender, input *TInput) (*Response[TOutput], error)
type StaticFunc[TOutput any] func(ctx context.Context, sender RequestSender) (*Response[TOutput], error)
type AuthenticatedRequestFunc[TInput any, TOutput any] func(ctx context.Context, sender RequestSender, input *TInput, token authentication.Token) (*Response[TOutput], error)

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

func Create[TInput any, TOutput any](info RequestInfo) RequestFunc[TInput, TOutput] {
	pattern := NewPattern(info.Path)

	err := validate[TInput](pattern)
	if err != nil {
		panic(err)
	}

	return func(bCtx context.Context, sender RequestSender, input *TInput) (*Response[TOutput], error) {
		// Split the input parameters into path, query, header, and body parameters.
		pathParameters, queryParameters, headerParameters, bodyParameters, err := GetRequestParameters(input)
		if err != nil {
			return nil, err
		}

		requestKey := createRequestKey(pathParameters, queryParameters, headerParameters)

		ctx := WithRequestContext(bCtx, info, requestKey, input)

		path, err := pattern.String(pathParameters)
		if err != nil {
			return nil, err
		}

		url := GetHost(ctx).JoinPath(path)
		url.RawQuery = queryParameters.Encode()

		req, err := http.NewRequestWithContext(ctx, info.Method, url.String(), bodyParameters)
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
			if errData, err := UnmarshalErrorJSON(data); err == nil {
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

func CreateAuthenticated[TInput any, TOutput any](req RequestInfo) AuthenticatedRequestFunc[TInput, TOutput] {
	base := Create[TInput, TOutput](req)

	return func(ctx context.Context, sender RequestSender, input *TInput, token authentication.Token) (*Response[TOutput], error) {
		ctx = authentication.WithToken(ctx, token)
		return base(ctx, sender, input)
	}
}

func CreateStatic[TOutput any](info RequestInfo) StaticFunc[TOutput] {
	base := Create[struct{}, TOutput](info)

	return func(ctx context.Context, sender RequestSender) (*Response[TOutput], error) {
		return base(ctx, sender, nil)
	}
}
