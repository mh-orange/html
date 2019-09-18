package scraper

import (
	"os"
	"testing"

	"github.com/go-test/deep"
)

type a struct {
	Name string `html:".name"`
}

type b struct {
	Names []a `html:".value"`
}

type bb struct {
	Names []*a `html:".value"`
}

type c struct {
	B *b `html:".set"`
}

func TestDecoder(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		options   []Option
		want      interface{}
		got       interface{}
		wantErr   error
	}{
		{"A", "testdata/a.html", nil, &a{"  Hello World  "}, &a{}, nil},
		{"A (trim)", "testdata/a.html", []Option{TrimSpace()}, &a{"Hello World"}, &a{}, nil},
		{"B", "testdata/b.html", nil, &b{[]a{{"one"}, {"two"}, {"three"}}}, &b{}, nil},
		{"BB", "testdata/b.html", nil, &bb{[]*a{{"one"}, {"two"}, {"three"}}}, &bb{}, nil},
		{"C", "testdata/c.html", nil, &c{&b{[]a{{"one"}, {"two"}, {"three"}}}}, &c{}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := os.Open(test.inputFile)
			if err != nil {
				t.Fatalf("Failed to open %s: %v", test.inputFile, err)
			}

			dec := NewDecoder(file, test.options...)
			gotErr := dec.Decode(test.got)
			if test.wantErr == gotErr {
				if gotErr == nil {
					if diff := deep.Equal(test.want, test.got); diff != nil {
						t.Error(diff)
					}
				}
			} else {
				t.Errorf("Wanted error %v got %v", test.wantErr, gotErr)
			}
		})
	}
}
