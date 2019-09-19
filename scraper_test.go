package scraper_test

import (
	"fmt"

	"github.com/mh-orange/scraper"
)

func ExampleUnmarshal() {
	// Parse and unmarshal an HTML document into a very basic Go struct
	document := `<html><body><h1 id="name">Hello Scraper!</h1><a href="https://github.org/mh-orange/scraper">Scraper</a> is Grrrrrreat!</body></html>`
	v := &struct {
		Name string `scraper:"#name"`                    // Match the element with the ID "name" and set Name to the text content
		URL  string `scraper:"a" scrapeType:"attr:href"` // match the first A element and set URL to the value of the HREF attribute
	}{}
	err := scraper.Unmarshal([]byte(document), v)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", v)
	// Output: &{Name:Hello Scraper! URL:https://github.org/mh-orange/scraper}
}

func ExampleUnmarshal_slice() {
	// Parse and unmarshal an HTML document into a Go struct
	// with a slice of things
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
		// Match the element with the ID "name" and set Name to the text content
		Name string `scraper:"#name"`

		// find all elements matching "ul li", for each of them, create a new entry in
		// Items with the value set to the text content of the matching nodes
		Items []string `scraper:"ul li"`
	}{}
	err := scraper.Unmarshal([]byte(document), v)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", v)
	// Output: &{Name:Hello Scraper! Items:[Item 1 Item 2 Item 3]}
}
