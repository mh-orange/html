package scraper

import (
	"io"

	"golang.org/x/net/html"
)

type Option func(*Unmarshaler) error

func TrimSpace() Option {
	return func(u *Unmarshaler) error {
		u.trimSpace = true
		return nil
	}
}

type Decoder struct {
	r       io.Reader
	options []Option
}

func NewDecoder(r io.Reader, options ...Option) *Decoder {
	dec := &Decoder{
		r:       r,
		options: options,
	}
	return dec
}

func (dec *Decoder) Decode(v interface{}) error {
	root, err := html.Parse(dec.r)
	if err == nil {
		err = NewUnmarshaler(root, dec.options...).Unmarshal(v)
	}

	return err
}
