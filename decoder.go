package html

import (
	"io"
	"reflect"
	"strings"

	"golang.org/x/net/html"
)

type DecoderOption func(*Decoder) error

func TrimSpace() DecoderOption {
	return func(dec *Decoder) error {
		dec.trimSpace = true
		return nil
	}
}

type Decoder struct {
	r         io.Reader
	trimSpace bool
	err       error
}

func NewDecoder(r io.Reader, options ...DecoderOption) *Decoder {
	dec := &Decoder{r: r}
	for _, option := range options {
		dec.err = option(dec)
		if dec.err != nil {
			break
		}
	}
	return dec
}

func (dec *Decoder) Decode(v interface{}) error {
	if dec.err != nil {
		return dec.err
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		dec.err = &InvalidUnmarshalError{reflect.TypeOf(v)}
		return dec.err
	}

	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Struct {
		dec.err = &InvalidUnmarshalError{reflect.TypeOf(v)}
		return dec.err
	}

	var root *html.Node
	root, dec.err = html.Parse(dec.r)
	if dec.err == nil {
		dec.decode(&selection{root}, rv)
	}

	return dec.err
}

func (dec *Decoder) decode(n *selection, rv reflect.Value) {
	rt := rv.Type()
	for i := 0; dec.err == nil && i < rt.NumField(); i++ {
		ft := rt.Field(i)
		var t *tag
		if t, dec.err = parseTag(ft); dec.err == nil {
			dec.walk(n, &field{rv.Field(i), t})
		} else if dec.err == ErrNoTag {
			dec.err = nil
		}
	}
}

func (dec *Decoder) walk(n *selection, f *field) {
	if n.Type == html.ElementNode {
		if f.tag.matches(n) {
			if f.Kind() == reflect.Slice {
				var value reflect.Value
				if f.Type().Elem().Kind() == reflect.Ptr {
					value = reflect.New(f.Type().Elem().Elem())
				} else {
					value = reflect.New(f.Type().Elem())
				}
				dec.set(&field{value, f.tag}, n)
				if dec.err == nil {
					if f.Type().Elem().Kind() == reflect.Ptr {
						f.Set(reflect.Append(f.Value, value))
					} else {
						f.Set(reflect.Append(f.Value, reflect.Indirect(value)))
					}
				}
			} else {
				dec.set(f, n)
			}
			// short circuit
			return
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		dec.walk(&selection{c}, f)
		if dec.err != nil {
			break
		}
	}
}

func (dec *Decoder) set(f *field, n *selection) {
	if dec.err = f.unmarshal(n); dec.err == nil || (dec.err != ErrNoUnmarshaler) {
		return
	}
	// reset error in case it was ErrNoUnmarshaler (see above conditional)
	dec.err = nil

	if f.Kind() == reflect.Ptr {
		if f.IsNil() {
			f.Set(reflect.New(f.Type().Elem()))
		}
		f.Value = reflect.Indirect(f.Value)
	}

	if f.Kind() == reflect.Struct {
		dec.decode(n, f.Value)
	} else {
		value := n.value(f.tag)
		if dec.trimSpace {
			value = strings.TrimSpace(value)
		}
		dec.err = f.set(value)
	}
}
