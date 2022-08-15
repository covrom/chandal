package gjsonb

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JsonB[T any] struct {
	Val *T
}

func (n *JsonB[T]) Scan(value interface{}) error {
	if value == nil {
		n.Val = nil
		return nil
	}
	switch val := value.(type) {
	case []byte:
		if len(val) == 0 || bytes.EqualFold(val, []byte("null")) {
			n.Val = nil
			return nil
		}

		n.Val = new(T)
		return json.Unmarshal(val, n.Val)
	}
	return fmt.Errorf("unsupported database data type %T, needs []byte", value)
}

func (n JsonB[T]) Value() (driver.Value, error) {
	if n.Val == nil {
		return "null", nil
	}
	j, err := json.Marshal(n.Val)
	return string(j), err
}

func (n JsonB[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Val)
}

func (n *JsonB[T]) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || bytes.EqualFold(b, []byte("null")) {
		n.Val = nil
		return nil
	}
	n.Val = new(T)
	return json.Unmarshal(b, n.Val)
}

func (n JsonB[T]) String() string {
	b, _ := json.Marshal(n.Val)
	return string(b)
}
