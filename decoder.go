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
	"io"

	"golang.org/x/net/html"
)

// Decoder will read from an io.Reader, parse the content
// into a root *html.Node and then unmarshal the content
// into a receiver
type Decoder struct {
	r       io.Reader
	options []Option
}

// NewDecoder initializes a decoder for the given reader and options
func NewDecoder(r io.Reader, options ...Option) *Decoder {
	dec := &Decoder{
		r:       r,
		options: options,
	}
	return dec
}

// Decode the input stream and unmarshal it into v
func (dec *Decoder) Decode(v interface{}) error {
	root, err := html.Parse(dec.r)
	if err == nil {
		err = NewUnmarshaler(root, dec.options...).Unmarshal(v)
	}

	return err
}
