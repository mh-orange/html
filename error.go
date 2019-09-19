package scraper

import (
	"reflect"
)

// An UnmarshalTypeError describes a value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value string       // description of JSON value - "bool", "array", "number -5"
	Type  reflect.Type // type of Go value it could not be assigned to
}

func (e *UnmarshalTypeError) Error() string {
	return "scraper: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
	Want reflect.Kind
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "scraper: Unmarshal(nil)"
	}

	if e.Type.Kind() != e.Want {
		return "scraper: Unmarshal(non-" + e.Want.String() + " " + e.Type.Kind().String() + ")"
	}

	return "scraper: Unmarshal(nil " + e.Type.String() + ")"
}
