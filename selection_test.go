package html

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestSelectionValue(t *testing.T) {
	tests := []struct {
		name      string
		inputHtml string
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
			node, err := html.ParseFragment(strings.NewReader(test.inputHtml), ctx)
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
