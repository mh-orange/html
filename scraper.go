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

// Package scraper provides a means to parse and unmarshal HTML into
// Go structs.  Usage is best described by example:
//
//		package main
//
//		import (
//			"fmt"
//
//			"github.com/mh-orange/scraper"
//		)
//
//		type MyType struct {
//			Name string `scraper:"#name"`
//			URL  string `scraper:"a" scrapeType:"attr:href"`
//		}
//
//		func main() {
//			document := `<html><body><h1 id="name">Hello Scraper!</h1><a href="https://github.org/mh-orange/scraper">Scraper</a> is Grrrrrreat!</body></html>`
//			v := &MyType{}
//			err := scraper.Unmarshal([]byte(document), v)
//			if err != nil {
//				panic(err.Error())
//			}
//			fmt.Printf("%+v\n", v)
//			// &{Name:Hello Scraper! URL:https://github.org/mh-orange/scraper}
//		}
//
// Structs are unmarshaled by matching CSS selectors to elements in an html document
// tree.  Scraper uses the wonderful Cascadia (https://github.com/andybalholm/cascadia)
// package to parse and match CSS selectors.
//
// To specify matching and unmarshaling rules, use the "scraper" and "scrapeType" struct
// field tags.  The "scraper" tag is used to define the CSS selector and the "scrapeType"
// indicates whether the value should be the text content or an attribute
// of the matching element.  The default type (if the scrapeTag is omitted) is to use
// the text content.  For example, to match an element with the id "name" and
// capture its text content:
//		type MyType struct {
//			Name string `scraper:"#name"`
//		}
//
// Another example, which uses the href attribute of a matching "a" tag:
//		type MyType struct {
//			URL string `scraper:"a" scrapeType:"attr:href"`
//		}
// Note that the attribute name is specified after the type (attr) and a separating colon.
//
package scraper
