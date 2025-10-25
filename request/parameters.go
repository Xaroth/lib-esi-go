package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

func getRequestBodyJSON(val any) (io.Reader, error) {
	buf, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}

func requestValue(val any) (string, error) {
	switch value := val.(type) {
	case string:
		return value, nil
	case int:
		return strconv.Itoa(value), nil
	case int64:
		return strconv.FormatInt(value, 10), nil
	default:
		return "", fmt.Errorf("%w: %T", ErrInvalidValueType, value)
	}
}

func GetRequestParameters[TInput any](input *TInput) (map[string]any, url.Values, http.Header, io.Reader, error) {
	typ := reflect.TypeOf(input).Elem()

	pathParameters := make(map[string]any, 0)
	queryParameters := make(url.Values, 0)
	headerParameters := make(http.Header, 0)
	var bodyParameters io.Reader = nil

	for i := range typ.NumField() {
		field := typ.Field(i)
		value := reflect.ValueOf(input).Elem().Field(i).Interface()

		if tag, ok := field.Tag.Lookup("path"); ok {
			pathParameters[tag] = value
		}
		if tag, ok := field.Tag.Lookup("query"); ok {
			val, err := requestValue(value)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			queryParameters.Add(tag, val)
		}
		if tag, ok := field.Tag.Lookup("header"); ok {
			val, err := requestValue(value)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			headerParameters.Add(tag, val)
		}
		if tag, ok := field.Tag.Lookup("body"); ok {
			if tag == "" {
				tag = "json"
			}

			switch tag {
			case "json":
				body, err := getRequestBodyJSON(value)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				headerParameters.Set("Content-Type", "application/json")
				bodyParameters = body
			default:
				return nil, nil, nil, nil, fmt.Errorf("%w: %s", ErrInvalidBodyType, tag)
			}
		}
	}
	return pathParameters, queryParameters, headerParameters, bodyParameters, nil
}
