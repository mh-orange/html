package scraper_test

import (
	"fmt"

	"github.com/mh-orange/scraper"
)

type MyType struct {
	Name string `scraper:"#name"`
	URL  string `scraper:"a" scrapeType:"attr:href"`
}

func ExampleUnmarshal() {
	document := `<html><body><h1 id="name">Hello Scraper!</h1><a href="https://github.org/mh-orange/scraper">Scraper</a> is Grrrrrreat!</body></html>`
	v := &MyType{}
	err := scraper.Unmarshal([]byte(document), v)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", v)
	// Output: &{Name:Hello Scraper! URL:https://github.org/mh-orange/scraper}
}
