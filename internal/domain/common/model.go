package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// JSONB represents a JSON field that can handle any JSON structure
type JSONB map[string]interface{}

// Value Marshal - implements driver.Valuer
func (a JSONB) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan Unmarshal - implements sql.Scanner
func (a *JSONB) Scan(value any) error {
	if value == nil {
		*a = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, a)
}

// StringSlice is a custom type for handling JSON arrays in database
type StringSlice []string

// Scan implements the sql.Scanner interface
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	switch v := value.(type) {
	case string:
		if v == "" || v == "[]" {
			*s = []string{}
			return nil
		}
		return json.Unmarshal([]byte(v), s)
	case []byte:
		if len(v) == 0 || string(v) == "[]" {
			*s = []string{}
			return nil
		}
		return json.Unmarshal(v, s)
	default:
		return fmt.Errorf("cannot scan %T into StringSlice", value)
	}
}

// Value implements the driver.Valuer interface
func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(jsonBytes), nil
}

// StringMap is a custom type for handling JSON maps in database
type StringMap map[string]string

// Scan implements the sql.Scanner interface
func (m *StringMap) Scan(value interface{}) error {
	if value == nil {
		*m = make(map[string]string)
		return nil
	}

	switch v := value.(type) {
	case string:
		if v == "" || v == "{}" {
			*m = make(map[string]string)
			return nil
		}
		return json.Unmarshal([]byte(v), m)
	case []byte:
		if len(v) == 0 || string(v) == "{}" {
			*m = make(map[string]string)
			return nil
		}
		return json.Unmarshal(v, m)
	default:
		return fmt.Errorf("cannot scan %T into StringMap", value)
	}
}

// Value implements the driver.Valuer interface
func (m StringMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(jsonBytes), nil
}
