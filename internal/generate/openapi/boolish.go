package openapi

import (
	"encoding/json"
	"fmt"
)

// Boolish accepts JSON true/false or "true"/"false" strings.
type Boolish bool

func (b *Boolish) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	switch data[0] {
	case 't', 'f':
		var v bool
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		*b = Boolish(v)
		return nil
	case '"':
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		switch s {
		case "true":
			*b = true
		case "false":
			*b = false
		default:
			return fmt.Errorf("boolish: invalid string %q", s)
		}
		return nil
	default:
		return fmt.Errorf("boolish: unexpected JSON %s", string(data))
	}
}

func (b Boolish) IsTrue() bool { return bool(b) }
