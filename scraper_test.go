package scraper_test

import (
	"fmt"
	"strings"

	"github.com/mh-orange/scraper"
)

func ExampleUnmarshal() {
	// Parse and unmarshal an HTML document into a very basic Go struct
	document := `<html><body><h1 id="name">Hello Scraper!</h1><a href="https://github.org/mh-orange/scraper">Scraper</a> is Grrrrrreat!</body></html>`
	v := &struct {
		// Name is assigned the text content from the element with the ID "name"
		Name string `scraper:"#name"`

		// URL is assigned the HREF attribute of the first A element found
		URL string `scraper:"a" scrapeType:"attr:href"`
	}{}
	err := scraper.Unmarshal([]byte(document), v)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", v)
	// Output: &{Name:Hello Scraper! URL:https://github.org/mh-orange/scraper}
}

func ExampleUnmarshal_slice() {
	// Scraper can be used to unmarshal structs with slices
	// of things as well
	document := `
		<html>
			<body>
				<h1 id="name">Hello Scraper!</h1>
				<ul>
					<li>Item 1</li>
					<li>Item 2</li>
					<li>Item 3</li>
				</ul>
			</body>
		</html>`
	v := &struct {
		// Name is assigned the text content from the element with the ID "name"
		Name string `scraper:"#name"`

		// Items is appended with the text content of each element matching the
		// "ul li" CSS selector
		Items []string `scraper:"ul li"`
	}{}
	err := scraper.Unmarshal([]byte(document), v)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", v)
	// Output: &{Name:Hello Scraper! Items:[Item 1 Item 2 Item 3]}
}

func ExampleUnmarshal_nested() {
	// Scraper can be used to unmarshal structs with other structs
	// in them
	document := `
		<html>
			<body>
				<h1 id="name">Hello Scraper!</h1>
				<ul>
					<li>Item 1</li>
					<li>Item 2</li>
					<li>Item 3</li>
				</ul>
			</body>
		</html>`
	v := &struct {
		// Name is assigned the text content from the element with the ID "name"
		Name string `scraper:"#name"`

		// Items is matched with the ul tag and then names is matched by the
		// li tags within.  Nested structs will be unmarshaled with the matching
		// _subtree_ not the entire document
		Items struct {
			Names []string `scraper:"li"`
		} `scraper:"ul"`
	}{}
	err := scraper.Unmarshal([]byte(document), v)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", v)
	// Output: &{Name:Hello Scraper! Items:{Names:[Item 1 Item 2 Item 3]}}
}

func ExampleDecoder() {
	// Decoder is useful for unmarshaling from an input stream
	document := `<html><body><h1 id="name">Hello Scraper!</h1></body></html>`
	v := &struct {
		// Name is assigned the text content from the element with the ID "name"
		Name string `scraper:"#name"`
	}{}

	reader := strings.NewReader(document)
	scraper.NewDecoder(reader).Decode(v)
	fmt.Printf("%+v\n", v)
	// Output: &{Name:Hello Scraper!}
}
