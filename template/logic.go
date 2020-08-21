package template

import (
	"html/template"
	"io/ioutil"

	"github.com/MihaiBlebea/dog-ceo/dog"
)

type service struct {
	dogService dog.Service
}

// New returns a new template service
func New(dogService dog.Service) Service {
	return &service{dogService}
}

func (s *service) Load(path string) (*Page, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tmp, err := template.New("Template").Parse(string(b))
	if err != nil {
		return nil, err
	}
	// fetch breeds and dogs
	dogs, err := s.dogService.AllDogs()
	if err != nil {
		return nil, err
	}

	return &Page{tmp, dogs}, nil
}
