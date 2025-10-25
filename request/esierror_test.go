package request_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/xaroth/lib-esi-go/request"
)

func TestErrorData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		data              string
		expectedErrorData *request.ErrorData
		expectedErr       error
	}{
		{
			name: "success: simple error",
			data: `{"error":"too few items for 'characters', 'characters' is required"}`,
			expectedErrorData: &request.ErrorData{
				ErrorMessage: "too few items for 'characters', 'characters' is required",
				Details:      nil,
			},
		},
		{
			name: "success: huma error",
			data: `{"error":"too few items for 'characters', 'characters' is required", "details": [{"message": "too few items for 'characters'", "location": "characters", "value": "characters"}]}`,
			expectedErrorData: &request.ErrorData{
				ErrorMessage: "too few items for 'characters', 'characters' is required",
				Details: []request.ErrorDetails{
					{
						Message:  "too few items for 'characters'",
						Location: "characters",
						Value:    "characters",
					},
				},
			},
		},
		{
			name:              "failure: invalid JSON",
			data:              `-`,
			expectedErrorData: nil,
			expectedErr:       &json.SyntaxError{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			errData, err := request.UnmarshalErrorJSON([]byte(testCase.data))

			if testCase.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if reflect.TypeOf(err) != reflect.TypeOf(testCase.expectedErr) {
					t.Fatalf("expected error %v, got %v", testCase.expectedErr, err)
				}
				if errData != nil {
					t.Fatalf("expected no error data, got: %+v", errData)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(testCase.expectedErrorData, errData); diff != "" {
				t.Fatalf("error data mismatch: %s", diff)
			}
		})
	}
}
