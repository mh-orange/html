package scraper_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mh-orange/scraper"
)

type Name struct {
	First string
	Last  string
}

func (n *Name) UnmarshalText(text []byte) (err error) {
	tokens := strings.Split(string(text), ", ")
	if len(tokens) == 2 {
		n.Last = tokens[0]
		n.First = tokens[1]
	} else {
		err = errors.New("Wanted comma separated last and first names")
	}
	return err
}

type Class struct {
	Students []Name `scraper:"ul li"`
}

func ExampleTextUnmarshaler() {
	document := `
		<html>
			<body>
				<h1 id="name">Class Roster</h1>
				<ul>
					<li>Stone, John</li>
					<li>Priya, Ponnappa</li>
					<li>Wong, Mia</li>
				</ul>
			</body>
		</html>`
	v := &Class{}
	err := scraper.Unmarshal([]byte(document), v)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", v)
	// Output: &{Students:[{First:John Last:Stone} {First:Ponnappa Last:Priya} {First:Mia Last:Wong}]}
}
