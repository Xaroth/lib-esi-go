package pattern

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrMissingVariable    = errors.New("missing variable")
	ErrInvalidVariable    = errors.New("invalid variable")
	ErrExtraneousVariable = errors.New("extraneous variable(s)")
	ErrUndefinedVariable  = errors.New("undefined variable")
)

type Pattern interface {
	String(variables map[string]any) (string, error)
	Variables() []string
	Validate(typ reflect.Type) error
}

type segment struct {
	variable string
	value    string
}

type pattern struct {
	segments []segment
}

// New creates a new pattern from the given path.
func New(path string) Pattern {
	segments := make([]segment, 0)

	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			segments = append(segments, segment{
				variable: part[1 : len(part)-1],
				value:    "",
			})
		} else {
			segments = append(segments, segment{
				variable: "",
				value:    part,
			})
		}
	}

	return &pattern{
		segments: segments,
	}
}

// NewValidated creates a new pattern and validates it against the given input type.
func NewValidated[TInput any](path string) (Pattern, error) {
	pattern := New(path)
	err := Validate[TInput](pattern)
	if err != nil {
		return nil, err
	}
	return pattern, nil
}

func formatVariable(variable any) (string, error) {
	switch value := variable.(type) {
	case string:
		return value, nil
	case int:
		return strconv.Itoa(value), nil
	case int64:
		return strconv.FormatInt(value, 10), nil
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	case fmt.Stringer:
		return value.String(), nil
	default:
		return "", fmt.Errorf("%w: %T", ErrInvalidVariable, value)
	}
}

func (p *pattern) Variables() []string {
	variables := make([]string, 0)
	for _, segment := range p.segments {
		if segment.variable != "" {
			variables = append(variables, segment.variable)
		}
	}
	return variables
}

func (p *pattern) String(variables map[string]any) (string, error) {
	parts := make([]string, len(p.segments))
	used := make(map[string]bool)

	for i, segment := range p.segments {
		if segment.variable != "" {
			variable, ok := variables[segment.variable]
			if !ok {
				return "", fmt.Errorf("%w: %s", ErrMissingVariable, segment.variable)
			}
			formatted, err := formatVariable(variable)
			if err != nil {
				return "", err
			}
			parts[i] = formatted
			used[segment.variable] = true
		} else {
			parts[i] = segment.value
		}
	}
	for variable := range variables {
		if _, ok := used[variable]; !ok {
			return "", fmt.Errorf("%w: %s", ErrExtraneousVariable, variable)
		}
	}

	return strings.Join(parts, "/"), nil
}

func (p *pattern) Validate(typ reflect.Type) error {
	variables := p.Variables()
	foundVariables := make(map[string]bool)

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

func Validate[T any](pattern Pattern) error {
	typ := reflect.TypeOf(new(T)).Elem()
	return pattern.Validate(typ)
}
