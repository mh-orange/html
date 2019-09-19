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
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/andybalholm/cascadia"
)

var (
	// ErrUnknownTagType indicates that the scraperType tag is an unknown value
	ErrUnknownTagType = errors.New("Unknown tag type ")

	errNoTag = errors.New("No HTML Tag found")
)

// Scraper uses struct field tags to determine how to unmarshal an HTML element tree into
// a type.  This is similar to how encoding/json uses tags to match json field names to struct
// field names.  There are two tags that scraper uses in its processing, `scraper` and `scrapeType`.
// Example:
//   type MyType struct {
//     URL string `scraper:"a.myurl" scrapeType:"attr:href"` // parses the href attribute from the matching a
//   }
const (
	// SelectorTagName is used to reflect the appropriate struct field tag.  The SelectorTagName
	// is the tag used to specify a CSS selector to match for the field
	SelectorTagName = "scraper"

	// TypeTagName (scrapeType) is the tag used to specify what kind of value lookup should be performed.  The
	// default is `text` and simply gathers the text nodes from the matching html subtree.  The
	// alternative type is `attr` which will assign value based on a matching attribute.  The
	// attribute name (for the matched node) is specified following a colon
	TypeTagName = "scrapeType"
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
	if tag, found := field.Tag.Lookup(SelectorTagName); found {
		err = t.parse(tag, field.Tag.Get(TypeTagName))
	} else {
		err = errNoTag
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
