package template

import (
	"html/template"
	"io"
	"time"

	"github.com/MihaiBlebea/dog-ceo/dog"
)

// Page struct
type Page struct {
	template *template.Template
	Dogs     []dog.Dog
	Duration time.Duration
}

// Render the page
func (p *Page) Render(w io.Writer, duration time.Duration) error {
	p.Duration = duration

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
