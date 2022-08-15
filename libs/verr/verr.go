package verr

import "encoding/json"

type JsonErr struct {
	Err error `json:"_err"`
}

type ValErr[T any] struct {
	Value *T
	Err   error
}

func (v ValErr[T]) MarshalJSON() ([]byte, error) {
	if v.Err != nil {
		return json.Marshal(JsonErr{Err: v.Err})
	}
	return json.Marshal(v.Value)
}
