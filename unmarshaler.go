// Copyright 2019 Andrew Bates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// Option updates an Unmarshaler with various capabilities
type Option func(*Unmarshaler) error

// TrimSpace tells the unmarshaller to trim values using strings.TrimSpace
// when a field is set, the value (either text content or attribute value)
// will be trimmed prior to type conversion and assignment
func TrimSpace() Option {
	return func(u *Unmarshaler) error {
		u.trimSpace = true
		return nil
	}
}

// BinaryUnmarshaler is the interface implemented by an object that can unmarshal
// the byte string (either text content or attribute) from an element matched
// by a scraper seleector
type BinaryUnmarshaler interface {
	encoding.BinaryUnmarshaler
}

// TextUnmarshaler is the interface implemented by an object that can unmarshal
// the byte string (either text content or attribute) from an element matched
// by a scraper seleector
type TextUnmarshaler interface {
	encoding.TextUnmarshaler
}

// HTMLUnmarshaler is the interface implemented by types that can unmarshal parsed
// html directly.  The input is a parsed element tree starting at the element that
// matched the CSS selector specified in the scraper tag
type HTMLUnmarshaler interface {
	UnmarshalHTML(*html.Node) error
}

// Unmarshal will parse the input text and unmarshal it into v
func Unmarshal(text []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(text)).Decode(v)
}

// Unmarshaler processes an HTML tree and unmarshals/parses it
// into a receiver.  The unmarshaler looks for struct field tags
// matching `scraper` and `scrapeType`
type Unmarshaler struct {
	root      *html.Node
	trimSpace bool
	err       error
}

// NewUnmarshaler creates a scraper Unmarshaler with its root set to the
// input *html.Node and setting any options given.  If any of the options
// generate an error, then that error is passed through upon calling
// Unmarshal.  This allows for chaining the NewUnmarshaler function with
// Unmarshal:
//   err := NewUnmarshaler(root).Unmarshal(v)
//
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

// Unmarshal the document into v
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
		} else if err == errNoTag {
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
