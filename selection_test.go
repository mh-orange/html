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
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestSelectionValue(t *testing.T) {
	tests := []struct {
		name      string
		inputHTML string
		inputTag  *tag
		want      string
	}{
		{"text", "<p>some value <strong>with emphasis</strong></p>", &tag{typ: text}, "some value with emphasis"},
		{"attr", `<a href="/some/path/somewhere">Click Me!</a>`, &tag{typ: attr, detail: "href"}, "/some/path/somewhere"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Div,
				Data:     "div",
			}
			node, err := html.ParseFragment(strings.NewReader(test.inputHTML), ctx)
			if err != nil {
				t.Fatalf("Failed to parse html: %v", err)
			}

			s := &selection{node[0]}
			got := s.value(test.inputTag)
			if test.want != got {
				t.Errorf("Wanted value %q got %q", test.want, got)
			}
		})
	}
}
