package html

import (
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

func TestTagParse(t *testing.T) {
	tests := []struct {
		name    string
		input   reflect.StructField
		wantTyp tagType
		wantDet string
		wantErr error
	}{
		{"A OK", reflect.StructField{Tag: `html:""`}, text, "", nil},
		{"A OK with detail", reflect.StructField{Tag: `html:"" htmlType:":foobar"`}, text, "foobar", nil},
		{"attribute instead of text", reflect.StructField{Tag: `html:"" htmlType:"attr"`}, attr, "", nil},
		{"unknown type", reflect.StructField{Tag: `html:"" htmlType:"foo"`}, text, "", ErrUnknownTagType},
		{"no tag", reflect.StructField{}, text, "", ErrNoTag},
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

type testMarshaler interface {
	called() bool
}

type testHtmlUnmarshaler struct{ c bool }

func (t testHtmlUnmarshaler) called() bool { return t.c }
func (t *testHtmlUnmarshaler) UnmarshalHtml(*html.Node) error {
	t.c = true
	return nil
}

type testTextUnmarshaler struct{ c bool }

func (t testTextUnmarshaler) called() bool { return t.c }
func (t *testTextUnmarshaler) UnmarshalText([]byte) error {
	t.c = true
	return nil
}

type testBinaryUnmarshaler struct{ c bool }

func (t testBinaryUnmarshaler) called() bool { return t.c }
func (t *testBinaryUnmarshaler) UnmarshalBinary([]byte) error {
	t.c = true
	return nil
}

type testNoUnmarshaler struct{ c bool }

func (t testNoUnmarshaler) called() bool { return t.c }

func TestFieldUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   testMarshaler
		want    bool
		wantErr error
	}{
		{"html", &testHtmlUnmarshaler{}, true, nil},
		{"text", &testTextUnmarshaler{}, true, nil},
		{"binary", &testBinaryUnmarshaler{}, true, nil},
		{"error", &testNoUnmarshaler{}, false, ErrNoUnmarshaler},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := &field{reflect.ValueOf(test.input), &tag{}}
			gotErr := f.unmarshal(&selection{})
			if test.wantErr == gotErr {
				if gotErr == nil {
					got := test.input.called()
					if test.want != got {
						t.Errorf("Wanted called to be %v got %v", test.want, got)
					}
				}
			} else {
				t.Errorf("Wanted error %v got %v", test.wantErr, gotErr)
			}
		})
	}
}
