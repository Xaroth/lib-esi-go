package enum_test

import (
	"testing"

	"github.com/xaroth/lib-esi-go/enum"
)

type TestType string
type TestEnum struct {
	Value1       TestType
	ValueWithTag TestType `value:"value_with_tag"`
}

var Test = enum.New[TestEnum]()
var Validator = enum.Validator[TestEnum, TestType]()

func TestNewEnum(t *testing.T) {
	t.Parallel()

	if Test.Value1 != "Value1" {
		t.Errorf("expected Value1 to be 'value1', got %s", Test.Value1)
	}
	if Test.ValueWithTag != "value_with_tag" {
		t.Errorf("expected ValueWithTag to be 'value_with_tag', got %s", Test.ValueWithTag)
	}
}

func TestEnumValidator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		value    TestType
		expected bool
	}{
		{
			name:     "success: enum value",
			value:    Test.Value1,
			expected: true,
		},
		{
			name:     "success: enum value with tag",
			value:    Test.ValueWithTag,
			expected: true,
		},
		{
			name:     "success: manual enum value",
			value:    "Value1",
			expected: true,
		},
		{
			name:     "failure: invalid value",
			value:    "invalid",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			valid := Validator(tc.value)
			if valid != tc.expected {
				t.Errorf("expected %t, got %t", tc.expected, valid)
			}
		})
	}
}
