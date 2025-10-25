package request

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrMissingVariable    = errors.New("missing variable")
	ErrInvalidVariable    = errors.New("invalid variable")
	ErrExtraneousVariable = errors.New("extraneous variable(s)")
)

type PatternVariable interface {
	PatternVariable() string
}

type Pattern interface {
	String(variables map[string]any) (string, error)
	Variables() []string
}

type segment struct {
	variable string
	value    string
}

type pattern struct {
	segments []segment
}

func NewPattern(path string) Pattern {
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
	case PatternVariable:
		return value.PatternVariable(), nil
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
