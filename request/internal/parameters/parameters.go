package parameters

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

var (
	ErrInvalidValueType = errors.New("invalid value type")
	ErrInvalidBodyType  = errors.New("invalid body type")
	ErrRequiredValue    = errors.New("required value is nil")
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
	case fmt.Stringer:
		return value.String(), nil
	default:
		return "", fmt.Errorf("%w: %T", ErrInvalidValueType, value)
	}
}

// Extract extracts the path, query, header, and body parameters from the input.
func Extract[TInput any](input *TInput) (map[string]any, url.Values, http.Header, io.Reader, error) {
	typ := reflect.TypeOf(input).Elem()

	pathParameters := make(map[string]any, 0)
	queryParameters := make(url.Values, 0)
	headerParameters := make(http.Header, 0)
	var bodyParameters io.Reader = nil

	for i := range typ.NumField() {
		field := typ.Field(i)
		fieldValue := reflect.ValueOf(input).Elem().Field(i)
		value := fieldValue.Interface()

		hasValue := fieldValue.IsValid()
		isZero := fieldValue.IsZero()

		isNil := !hasValue
		if !hasValue {
			isNil = fieldValue.IsNil()
		}

		if tag, ok := field.Tag.Lookup("path"); ok {
			if isZero || isNil {
				// Path parameters are always required.
				return nil, nil, nil, nil, ErrRequiredValue
			}
			pathParameters[tag] = value
		}

		// For all non-path parameters, if the value is nil, and the field is required
		// error early.
		if isZero || isNil {
			if tag, ok := field.Tag.Lookup("required"); ok && tag == "true" {
				return nil, nil, nil, nil, ErrRequiredValue
			}
			continue
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
