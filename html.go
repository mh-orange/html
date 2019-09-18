package html

/*type Html interface {
	find(selector string) (Html, error)
	each(selector string, callback func(Html) error) error
	text() string
	attr(attribute string) (string, bool)
}

type html struct {
	*goquery.Selection
}

func (h *html) attr(attribute string) (string, bool) {
	return h.Selection.Attr(attribute)
}

func (h *html) text() string {
	return h.Selection.Text()
}

func (h *html) find(selector string) (Html, error) {
	var err error
	if selector == "" {
		return h, nil
	}

	node := &html{}
	node.Selection = h.Selection.Find(selector)
	if node.Length() == 0 {
		err = &Error{Msg: fmt.Sprintf("Could not find %q", selector)}
	} else if node.Length() > 1 {
		err = &Error{Msg: fmt.Sprintf("%q returned multiple nodes", selector)}
	}
	return node, err
}

func (h *html) each(selector string, callback func(Html) error) (err error) {
	selection := h.Selection.Find(selector)
	selection.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		err = callback(&html{Selection: selection})
		return err == nil
	})
	return
}

func Decode(text string) (Html, error) {
	h := &html{}
	return h, h.UnmarshalText([]byte(text))
}

func (h *html) UnmarshalText(text []byte) (err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(text))
	if err == nil {
		h.Selection = doc.Selection
	}
	return
}*/
