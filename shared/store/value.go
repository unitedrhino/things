package store

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Int64Arr []int64

func (p Int64Arr) Value() (driver.Value, error) {
	if p == nil {
		return "[]", nil
	}
	valByte, _ := json.Marshal(p)
	val := string(valByte)
	if val == "" {
		val = "[]"
	}
	return val, nil
}
func (p *Int64Arr) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("failed to scan point: value is nil")
	}
	switch value.(type) {
	case []byte:
		va := value.([]byte)
		return json.Unmarshal(va, p)
	case string:
		va := value.(string)
		return json.Unmarshal([]byte(va), p)
	default:
		return fmt.Errorf("failed to scan point: invalid type: %T", value)
	}
	return nil
}
