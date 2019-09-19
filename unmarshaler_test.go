package scraper

import (
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

type testMarshaler interface {
	called() bool
}

type testHTMLUnmarshaler struct{ c bool }

func (t testHTMLUnmarshaler) called() bool { return t.c }
func (t *testHTMLUnmarshaler) UnmarshalHTML(*html.Node) error {
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

func TestTryUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   testMarshaler
		want    bool
		wantErr error
	}{
		{"html", &testHTMLUnmarshaler{}, true, nil},
		{"text", &testTextUnmarshaler{}, true, nil},
		{"binary", &testBinaryUnmarshaler{}, true, nil},
		{"error 1", &testNoUnmarshaler{}, false, errNoUnmarshaler},
		{"error 2", testNoUnmarshaler{}, false, errNoUnmarshaler},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u := &Unmarshaler{}
			f := &field{reflect.ValueOf(test.input), &tag{}}
			gotErr := u.tryUnmarshaler(f, &selection{})
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
