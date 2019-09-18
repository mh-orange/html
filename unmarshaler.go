package html

import (
	"encoding"
	"errors"

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

/*func Unmarshal(node *html.Node, v interface{}) error {
	un := &Unmarshaler{
		TrimSpace: true,
		Node:      node,
	}
	return un.Unmarshal(v)
}

type Unmarshaler struct {
	TrimSpace bool
	Node      *html.Node
}

func (un *Unmarshaler) Unmarshal(v interface{}) (err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Struct {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	return un.unmarshal(un.Html, rv)
}

func (un Unmarshaler) unmarshal(selection Html, rv reflect.Value) (err error) {
	rt := rv.Type()
	for i := 0; err == nil && i < rt.NumField(); i++ {
		ft := rt.Field(i)
		var t *tag
		if t, err = parseTag(ft); err == nil {
			if rv.Field(i).Kind() == reflect.Slice {
				err = selection.each(ti.selector, func(node Html) error {
					value := reflect.New(ft.Type.Elem())
					err := un.set(ti, node, value)
					if err == nil {
						if ft.Type.Elem().Kind() == reflect.Ptr {
							rv.Field(i).Set(reflect.Append(rv.Field(i), value))
						} else {
							rv.Field(i).Set(reflect.Append(rv.Field(i), reflect.Indirect(value)))
						}
					}
					return err
				})
			} else {
				var node Html
				node, err = selection.find(ti.selector)
				if err == nil {
					err = un.set(ti, node, rv.Field(i))
				}
			}
		} else if err == ErrNoTag {
			err = nil
		}
	}
	return err
}

func (un Unmarshaler) unmarshalField(ti *tagInfo, selection Html, field reflect.Value) (err error) {
	if field.Kind() != reflect.Slice && field.Kind() != reflect.Ptr {
		field = field.Addr()
	}
	if field.Type().NumMethod() > 0 && field.CanInterface() {
		switch u := field.Interface().(type) {
		case HtmlUnmarshaler:
			err = u.UnmarshalHtml(selection)
		case TextUnmarshaler:
			err = u.UnmarshalText([]byte(ti.value(selection, un.TrimSpace)))
		case BinaryUnmarshaler:
			err = u.UnmarshalBinary([]byte(ti.value(selection, un.TrimSpace)))
		}
	} else {
		err = ErrNoUnmarshaler
	}
	return err
}

func (un Unmarshaler) set(ti *tagInfo, selection Html, field reflect.Value) (err error) {
	if err := un.unmarshalField(ti, selection, field); err == nil || (err != ErrNoUnmarshaler) {
		return err
	}

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type()))
		}
		field = reflect.Indirect(field)
	}

	if field.Kind() == reflect.Struct {
		err = un.unmarshal(selection, field)
	} else {
		err = ti.set(selection, un.TrimSpace, field)
	}
	return err
}*/
