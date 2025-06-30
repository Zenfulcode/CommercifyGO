package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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
