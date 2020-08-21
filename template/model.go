package template

import (
	"html/template"
	"io"

	"github.com/MihaiBlebea/dog-ceo/dog"
)

// Page struct
type Page struct {
	template *template.Template
	Dogs     []dog.Dog
}

// Render the page
func (p *Page) Render(w io.Writer) error {
	err := p.template.Execute(w, p)
	if err != nil {
		return err
	}

	return nil
}

// DogsCount get dogs number by count
func (p *Page) DogsCount(count int) []dog.Dog {
	if len(p.Dogs) <= count {
		return p.Dogs
	}

	return p.Dogs[0:count]
}
