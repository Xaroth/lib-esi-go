package request

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"slices"

	"github.com/xaroth/lib-esi-go/util/esierror"
	"github.com/xaroth/lib-esi-go/util/pattern"
)

var (
	ErrInvalidValueType  = errors.New("invalid value type")
	ErrInvalidBodyType   = errors.New("invalid body type")
	ErrInvalidStatusCode = errors.New("invalid status code")

	ErrUndefinedVariable  = errors.New("undefined variable")
	ErrExtraneousVariable = errors.New("extraneous variable")

	ErrNotModified = errors.New("not modified")
)

// The bare minimum interface that we need to send a request.
type RequestSender interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestScheduler interface {
	Schedule(ctx context.Context, info RequestInfo, token Token) (func(*http.Response), error)
}

type SchedulingSender interface {
	RequestSender
	RequestScheduler
}

type Response[TOutput any] struct {
	*http.Response
	// If an error occurs, the Data field is not guaranteed to be set.
	Data TOutput
	// This field is only set if the status code is not in the 200-299 range.
	ErrorData *esierror.ErrorData
}

type RequestInfo struct {
	Method        string
	Path          string
	RequiredScope string
}

func validate[TInput any](pattern pattern.Pattern) error {
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

type RequestFunc[TInput any, TOutput any] func(ctx context.Context, client RequestSender, input *TInput, opts ...Option) (*Response[TOutput], error)
type StaticFunc[TOutput any] func(ctx context.Context, client RequestSender, opts ...Option) (*Response[TOutput], error)
type AuthenticatedRequestFunc[TInput any, TOutput any] func(ctx context.Context, client RequestSender, input *TInput, token Token, opts ...Option) (*Response[TOutput], error)

func Create[TInput any, TOutput any](info RequestInfo) RequestFunc[TInput, TOutput] {
	pattern := pattern.New(info.Path)

	err := validate[TInput](pattern)
	if err != nil {
		panic(err)
	}

	return func(bCtx context.Context, client RequestSender, input *TInput, opts ...Option) (*Response[TOutput], error) {
		// Split the input parameters into path, query, header, and body parameters.
		pathParameters, queryParameters, headerParameters, bodyParameters, err := GetRequestParameters(input)
		if err != nil {
			return nil, err
		}

		path, err := pattern.String(pathParameters)
		if err != nil {
			return nil, err
		}

		cfg := NewConfig(client, opts...)

		var notifyDone func(*http.Response) = func(*http.Response) {}
		var notifyError func() = func() {}

		// Allow the ESI Client to delay sending the request for rate-limiting purposes.
		// This means that we might need to refresh the token if it is expired, if the token is refreshable.
		if c, ok := client.(SchedulingSender); ok {
			ctx, cancel := context.WithTimeout(bCtx, cfg.scheduleTimeout)
			defer cancel()

			notifyDone, err = c.Schedule(ctx, info, cfg.token)
			notifyError = func() { notifyDone(nil) }
			if err != nil {
				return nil, err
			}

			if cfg.token != nil {
				if refreshable, ok := cfg.token.(RefreshableToken); ok {
					if err := refreshable.RefreshIfNeeded(ctx); err != nil {
						notifyError()
						return nil, err
					}
				}
			}
		}

		ctx, cancel := context.WithTimeout(bCtx, cfg.requestTimeout)
		defer cancel()

		url := cfg.Host().JoinPath(path)
		url.RawQuery = queryParameters.Encode()

		req, err := http.NewRequestWithContext(ctx, info.Method, url.String(), bodyParameters)
		if err != nil {
			notifyError()
			return nil, err
		}

		for key, values := range headerParameters {
			req.Header.Del(key)
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		cfg.ApplyRequest(req)

		rawResp, err := client.Do(req)
		if err != nil {
			notifyError()
			return nil, err
		}
		defer notifyDone(rawResp)
		defer rawResp.Body.Close()

		data, err := io.ReadAll(rawResp.Body)
		if err != nil {
			return nil, err
		}

		resp := &Response[TOutput]{
			Response: rawResp,
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			if errData, err := esierror.UnmarshalJSON(data); err == nil {
				resp.ErrorData = errData

			}

			switch resp.StatusCode {
			case http.StatusNotModified:
				return resp, ErrNotModified
			default:
				if resp.ErrorData != nil {
					return resp, fmt.Errorf("%w: %d: %s", ErrInvalidStatusCode, resp.StatusCode, resp.ErrorData.Error())
				}
				return resp, fmt.Errorf("%w: %d", ErrInvalidStatusCode, resp.StatusCode)
			}
		}

		if len(data) > 0 {
			if err := json.Unmarshal(data, &resp.Data); err != nil {
				return resp, err
			}
		}

		return resp, nil
	}
}

func CreateAuthenticated[TInput any, TOutput any](req RequestInfo) AuthenticatedRequestFunc[TInput, TOutput] {
	baseRequest := Create[TInput, TOutput](req)

	return func(ctx context.Context, client RequestSender, input *TInput, token Token, opts ...Option) (*Response[TOutput], error) {
		downstreamOptions := make([]Option, len(opts)+1)
		copy(downstreamOptions, opts)
		downstreamOptions = append(downstreamOptions, WithToken(token))

		return baseRequest(ctx, client, input, downstreamOptions...)
	}
}

func CreateStatic[TOutput any](info RequestInfo) StaticFunc[TOutput] {
	baseRequest := Create[struct{}, TOutput](info)

	return func(ctx context.Context, client RequestSender, opts ...Option) (*Response[TOutput], error) {
		return baseRequest(ctx, client, nil, opts...)
	}
}
