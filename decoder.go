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
