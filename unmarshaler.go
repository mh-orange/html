package scraper

import (
	"bytes"
	"encoding"
	"errors"
	"reflect"
	"strings"

	"golang.org/x/net/html"
)

var (
	errNoUnmarshaler = errors.New("Type does not implement a known umarshaler")
)

type BinaryUnmarshaler interface {
	encoding.BinaryUnmarshaler
}

type TextUnmarshaler interface {
	encoding.TextUnmarshaler
}

type HTMLUnmarshaler interface {
	UnmarshalHTML(*html.Node) error
}

func Unmarshal(text []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(text)).Decode(v)
}

type Unmarshaler struct {
	root      *html.Node
	trimSpace bool
	err       error
}

func NewUnmarshaler(root *html.Node, options ...Option) (u *Unmarshaler) {
	u = &Unmarshaler{root: root}
	for _, option := range options {
		u.err = option(u)
		if u.err != nil {
			break
		}
	}
	return u
}

func (u *Unmarshaler) Unmarshal(v interface{}) (err error) {
	if u.err != nil {
		return u.err
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v), reflect.Ptr}
	}

	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Struct {
		return &InvalidUnmarshalError{rv.Type(), reflect.Struct}
	}

	return u.unmarshalStruct(&field{Value: rv, tag: &tag{typ: text}}, &selection{u.root})
}

func (u *Unmarshaler) tryUnmarshaler(f *field, n *selection) error {
	value := f.Value
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Ptr {
		if value.CanAddr() {
			value = value.Addr()
		} else {
			return errNoUnmarshaler
		}
	}

	err := errNoUnmarshaler
	if value.Type().NumMethod() > 0 && value.CanInterface() {
		switch i := value.Interface().(type) {
		case TextUnmarshaler:
			err = i.UnmarshalText([]byte(u.value(f, n)))
		case BinaryUnmarshaler:
			err = i.UnmarshalBinary([]byte(u.value(f, n)))
		case HTMLUnmarshaler:
			err = i.UnmarshalHTML(n.Node)
		}
	}
	return err
}

// unmarshalStruct a struct
func (u *Unmarshaler) unmarshalStruct(f *field, n *selection) (err error) {
	if err = u.tryUnmarshaler(f, n); err != errNoUnmarshaler {
		return err
	}
	err = nil

	rt := f.Value.Type()
	for i := 0; err == nil && i < rt.NumField(); i++ {
		ft := rt.Field(i)
		var t *tag
		if t, err = parseTag(ft); err == nil {
			err = u.walk(&field{f.Value.Field(i), t}, n)
		} else if err == ErrNoTag {
			err = nil
		}
	}
	return
}

func (u *Unmarshaler) walk(f *field, n *selection) (err error) {
	if n.Type == html.ElementNode {
		if f.tag.matches(n) {
			// short circuit
			return u.unmarshalField(f, n)
		}
	}

	for c := n.FirstChild; err == nil && c != nil; c = c.NextSibling {
		err = u.walk(f, &selection{c})
	}
	return
}

func (u *Unmarshaler) value(f *field, n *selection) string {
	value := n.value(f.tag)
	if u.trimSpace {
		value = strings.TrimSpace(value)
	}
	return value
}

func (u *Unmarshaler) unmarshalField(f *field, n *selection) (err error) {
	if err = u.tryUnmarshaler(f, n); err != errNoUnmarshaler {
		return err
	}

	switch f.Kind() {
	case reflect.Slice:
		newField := &field{reflect.New(f.Type().Elem()), f.tag}
		err = u.unmarshalField(newField, n)
		if err == nil {
			if f.Type().Elem().Kind() == reflect.Ptr {
				f.Set(reflect.Append(f.Value, newField.Value.Addr()))
			} else {
				f.Set(reflect.Append(f.Value, reflect.Indirect(newField.Value)))
			}
		}
	case reflect.Struct:
		err = u.unmarshalStruct(f, n)
	case reflect.Ptr:
		if f.IsNil() {
			f.Set(reflect.New(f.Type().Elem()))
		}
		f.Value = reflect.Indirect(f.Value)
		err = u.unmarshalField(f, n)
	default:
		err = f.set(u.value(f, n))
	}

	return
}
