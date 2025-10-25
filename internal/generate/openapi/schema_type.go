package openapi

import (
	"encoding/json"
	"fmt"
)

// SchemaType accepts JSON string or string array (OpenAPI 3.1 nullable type).
type SchemaType string

func (t *SchemaType) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	switch data[0] {
	case '"':
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*t = SchemaType(s)
		return nil
	case '[':
		var types []string
		if err := json.Unmarshal(data, &types); err != nil {
			return err
		}
		for _, s := range types {
			if s != "null" {
				*t = SchemaType(s)
				return nil
			}
		}
		if len(types) > 0 {
			*t = SchemaType(types[0])
		}
		return nil
	default:
		return fmt.Errorf("schema type: unexpected JSON %s", string(data))
	}
}

func (t SchemaType) String() string { return string(t) }
