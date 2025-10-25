package request

import (
	"encoding/json"
	"fmt"
)

type ErrorDetails struct {
	Message  string `json:"message"`
	Location string `json:"location"`
	Value    string `json:"value"`
}

func (e ErrorDetails) Error() string {
	if e.Location != "" {
		return fmt.Sprintf("%s: %s: %s", e.Location, e.Message, e.Value)
	}
	return fmt.Sprintf("%s: %s", e.Message, e.Value)
}

type ErrorData struct {
	ErrorMessage string         `json:"error"`
	Details      []ErrorDetails `json:"details"`
}

func UnmarshalErrorJSON(data []byte) (*ErrorData, error) {
	var errData ErrorData
	if err := json.Unmarshal(data, &errData); err != nil {
		return nil, err
	}
	return &errData, nil
}

func (e ErrorData) Error() string {
	return e.ErrorMessage
}

func (e ErrorData) String() string {
	return e.ErrorMessage
}
