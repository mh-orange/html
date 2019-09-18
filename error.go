package scraper

import (
	"fmt"
	"reflect"
)

type Error struct {
	Msg   string
	Cause string
}

func (err *Error) Error() string {
	if err.Cause == "" {
		return err.Msg
	}
	return fmt.Sprintf("%s: %s", err.Msg, err.Cause)
}

// An UnmarshalTypeError describes a JSON value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value string       // description of JSON value - "bool", "array", "number -5"
	Type  reflect.Type // type of Go value it could not be assigned to
}

func (e *UnmarshalTypeError) Error() string {
	return "html: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "html: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "html: Unmarshal(non-pointer " + e.Type.String() + ")"
	}

	if e.Type.Elem().Kind() != reflect.Struct {
		return "html: Unmarshal(non-struct " + e.Type.String() + ")"
	}

	return "html: Unmarshal(nil " + e.Type.String() + ")"
}
