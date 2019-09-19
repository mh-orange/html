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
	"reflect"
	"testing"
)

func TestTagParse(t *testing.T) {
	tests := []struct {
		name    string
		input   reflect.StructField
		wantTyp tagType
		wantDet string
		wantErr error
	}{
		{"A OK", reflect.StructField{Tag: `scraper:""`}, text, "", nil},
		{"A OK with detail", reflect.StructField{Tag: `scraper:"" scrapeType:":foobar"`}, text, "foobar", nil},
		{"attribute instead of text", reflect.StructField{Tag: `scraper:"" scrapeType:"attr"`}, attr, "", nil},
		{"unknown type", reflect.StructField{Tag: `scraper:"" scrapeType:"foo"`}, text, "", ErrUnknownTagType},
		{"no tag", reflect.StructField{}, text, "", errNoTag},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tag, gotErr := parseTag(test.input)
			if test.wantErr == gotErr {
				if gotErr == nil {
					if test.wantTyp != tag.typ {
						t.Errorf("Wanted type %d got %d", test.wantTyp, tag.typ)
					}
					if test.wantDet != tag.detail {
						t.Errorf("wanted detail %q got %q", test.wantDet, tag.detail)
					}
				}
			} else {
				t.Errorf("Wanted error %v got %v", test.wantErr, gotErr)
			}
		})
	}
}

func TestTagMatches(t *testing.T) {
	tests := []struct {
		name  string
		input *selection
		tag   *tag
		want  bool
	}{
		{"nil selector", nil, &tag{}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.tag.matches(test.input)
			if test.want != got {
				t.Errorf("Wanted %v match got %v", test.want, got)
			}
		})
	}
}

func TestFieldSet(t *testing.T) {
	isErr := func(wantErr *UnmarshalTypeError, gotErr error) bool {
		if wantErr == nil && gotErr == nil {
			return true
		} else if wantErr == nil || gotErr == nil {
			return false
		}

		if ge, ok := gotErr.(*UnmarshalTypeError); ok {
			if *wantErr == *ge {
				return true
			}
		}
		return false
	}

	tests := []struct {
		name     string
		input    string
		receiver reflect.Value
		want     reflect.Value
		wantErr  *UnmarshalTypeError
	}{
		{"string", "foo", reflect.New(reflect.TypeOf("")), reflect.ValueOf("foo"), nil},
		{"int", "1234", reflect.New(reflect.TypeOf(1)), reflect.ValueOf(1234), nil},
		{"int error", "i1234", reflect.New(reflect.TypeOf(1)), reflect.Value{}, &UnmarshalTypeError{Value: "number " + "i1234", Type: reflect.TypeOf(1)}},
		{"uint", "5678", reflect.New(reflect.TypeOf(uint(1))), reflect.ValueOf(uint(5678)), nil},
		{"uint error", "i5678", reflect.New(reflect.TypeOf(uint(1))), reflect.Value{}, &UnmarshalTypeError{Value: "number " + "i5678", Type: reflect.TypeOf(uint(1))}},
		{"float", "9.1011", reflect.New(reflect.TypeOf(float32(1))), reflect.ValueOf(float32(9.1011)), nil},
		{"float error", "i9.1011", reflect.New(reflect.TypeOf(float32(1))), reflect.Value{}, &UnmarshalTypeError{Value: "number " + "i9.1011", Type: reflect.TypeOf(float32(1))}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := reflect.Indirect(test.receiver)
			f := &field{Value: got}
			gotErr := f.set(test.input)
			if isErr(test.wantErr, gotErr) {
				if gotErr == nil {
					if !reflect.DeepEqual(test.want.Interface(), got.Interface()) {
						t.Errorf("Wanted value %v got %v", test.want.Interface(), got.Interface())
					}
				}
			} else {
				t.Errorf("Wanted error %v got %v", test.wantErr, gotErr)
			}
		})
	}
}
