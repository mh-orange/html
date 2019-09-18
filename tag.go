package scraper

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/andybalholm/cascadia"
)

var (
	ErrUnknownTagType = errors.New("Unknown tag type ")
	ErrNoTag          = errors.New("No HTML Tag found")
)

type tagType int

func (tt *tagType) UnmarshalString(str string) (err error) {
	switch str {
	case "":
		fallthrough
	case "text":
		*tt = text
	case "attr":
		*tt = attr
	default:
		err = ErrUnknownTagType
	}
	return err
}

const (
	text tagType = iota
	attr
)

type tag struct {
	selector cascadia.Selector
	typ      tagType
	detail   string
}

func parseTag(field reflect.StructField) (t *tag, err error) {
	t = &tag{}
	if tag, found := field.Tag.Lookup("html"); found {
		err = t.parse(tag, field.Tag.Get("htmlType"))
	} else {
		err = ErrNoTag
	}
	return t, err
}

func (t *tag) matches(node *selection) bool {
	if t.selector == nil {
		return true
	}
	return t.selector.Match(node.Node)
}

func (t *tag) parse(tagStr, typeStr string) (err error) {
	if tagStr != "" {
		t.selector, err = cascadia.Compile(tagStr)
	}
	if err == nil {
		typFields := strings.Split(typeStr, ":")
		err = t.typ.UnmarshalString(typFields[0])
		if len(typFields) > 1 {
			t.detail = typFields[1]
		}
	}
	return err
}

type field struct {
	reflect.Value
	tag *tag
}

func (f *field) set(value string) error {
	switch f.Kind() {
	case reflect.String:
		f.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil || f.OverflowInt(n) {
			return &UnmarshalTypeError{Value: "number " + value, Type: f.Type()}
		}
		f.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil || f.OverflowUint(n) {
			return &UnmarshalTypeError{Value: "number " + value, Type: f.Type()}
		}
		f.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(value, f.Type().Bits())
		if err != nil || f.OverflowFloat(n) {
			return &UnmarshalTypeError{Value: "number " + value, Type: f.Type()}
		}
		f.SetFloat(n)
	}
	return nil
}

func (f *field) unmarshal(n *selection) (err error) {
	value := f.Value

	if value.Kind() != reflect.Slice && value.Kind() != reflect.Ptr {
		if value.CanAddr() {
			value = value.Addr()
		} else {
			return ErrNoUnmarshaler
		}
	}

	if value.Type().NumMethod() > 0 && value.CanInterface() {
		switch u := value.Interface().(type) {
		case HtmlUnmarshaler:
			err = u.UnmarshalHtml(n.Node)
		case TextUnmarshaler:
			err = u.UnmarshalText([]byte(n.value(f.tag)))
		case BinaryUnmarshaler:
			err = u.UnmarshalBinary([]byte(n.value(f.tag)))
		}
	} else {
		err = ErrNoUnmarshaler
	}
	return err
}
