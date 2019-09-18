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
	ErrNoUnmarshaler = errors.New("Type does not implement a known umarshaler")
)

type BinaryUnmarshaler interface {
	encoding.BinaryUnmarshaler
}

type TextUnmarshaler interface {
	encoding.TextUnmarshaler
}

type HtmlUnmarshaler interface {
	UnmarshalHtml(*html.Node) error
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

func (u *Unmarshaler) Unmarshal(v interface{}) error {
	if u.err != nil {
		return u.err
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Struct {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	return u.unmarshal(&selection{u.root}, rv)
}

func (u *Unmarshaler) unmarshal(n *selection, rv reflect.Value) (err error) {
	rt := rv.Type()
	for i := 0; err == nil && i < rt.NumField(); i++ {
		ft := rt.Field(i)
		var t *tag
		if t, err = parseTag(ft); err == nil {
			err = u.walk(n, &field{rv.Field(i), t})
		} else if err == ErrNoTag {
			err = nil
		}
	}
	return
}

func (u *Unmarshaler) walk(n *selection, f *field) (err error) {
	if n.Type == html.ElementNode {
		if f.tag.matches(n) {
			if f.Kind() == reflect.Slice {
				var value reflect.Value
				if f.Type().Elem().Kind() == reflect.Ptr {
					value = reflect.New(f.Type().Elem().Elem())
				} else {
					value = reflect.New(f.Type().Elem())
				}
				err = u.set(&field{value, f.tag}, n)
				if err == nil {
					if f.Type().Elem().Kind() == reflect.Ptr {
						f.Set(reflect.Append(f.Value, value))
					} else {
						f.Set(reflect.Append(f.Value, reflect.Indirect(value)))
					}
				}
			} else {
				err = u.set(f, n)
			}
			// short circuit
			return
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err = u.walk(&selection{c}, f)
		if err != nil {
			break
		}
	}
	return
}

func (u *Unmarshaler) set(f *field, n *selection) (err error) {
	if err = f.unmarshal(n); err == nil || (err != ErrNoUnmarshaler) {
		return
	}
	// reset error in case it was ErrNoUnmarshaler (see above conditional)
	err = nil

	if f.Kind() == reflect.Ptr {
		if f.IsNil() {
			f.Set(reflect.New(f.Type().Elem()))
		}
		f.Value = reflect.Indirect(f.Value)
	}

	if f.Kind() == reflect.Struct {
		err = u.unmarshal(n, f.Value)
	} else {
		value := n.value(f.tag)
		if u.trimSpace {
			value = strings.TrimSpace(value)
		}
		err = f.set(value)
	}
	return
}
