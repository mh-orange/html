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

type testErrStruct struct{}

func TestError(t *testing.T) {
	var tes *testErrStruct

	tests := []struct {
		name  string
		input error
		want  string
	}{
		{"UnmarshalTypeError", &UnmarshalTypeError{"foo", reflect.TypeOf(0)}, "scraper: cannot unmarshal foo into Go value of type int"},
		{"nil InvalidUnmarshalError", &InvalidUnmarshalError{reflect.TypeOf(nil), reflect.Ptr}, "scraper: Unmarshal(nil)"},
		{"int InvalidUnmarshalError", &InvalidUnmarshalError{reflect.TypeOf(0), reflect.Struct}, "scraper: Unmarshal(non-struct int)"},
		{"nil ptr InvalidUnmarshalError", &InvalidUnmarshalError{reflect.TypeOf(tes), reflect.Ptr}, "scraper: Unmarshal(nil *scraper.testErrStruct)"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.Error()
			if test.want != got {
				t.Errorf("Wanted error string %q got %q", test.want, got)
			}
		})
	}
}
