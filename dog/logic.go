package dog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/sirupsen/logrus"
)

const url = "https://dog.ceo/api"

type service struct {
	url    string
	logger *logrus.Logger
}

// New retruns a new dog ceo service
func New(logger *logrus.Logger) Service {
	return &service{url, logger}
}

func (s *service) breeds() ([]Breed, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/breeds/list/all",
			s.url,
		),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []Breed{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var data struct {
		Message map[string][]string `json:"message"`
		Status  string              `json:"status"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return []Breed{}, err
	}

	var breeds []Breed
	for breed := range data.Message {
		breeds = append(breeds, toBreed(breed))
	}

	sort.SliceStable(breeds, func(i, j int) bool {
		return breeds[i] < breeds[j]
	})

	return breeds, nil
}

func (s *service) dogs(breed Breed) ([]Dog, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/breed/%v/images",
			s.url,
			breed,
		),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []Dog{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var data struct {
		Message []string `json:"message"`
		Status  string   `json:"status"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return []Dog{}, err
	}

	var dogs []Dog
	for _, dog := range data.Message {
		dogs = append(dogs, toDog(dog))
	}

	return dogs, nil
}

func (s *service) AllDogs() ([]Dog, error) {
	breeds, err := s.breeds()
	if err != nil {
		return nil, err
	}

	type result struct {
		index int
		dogs  []Dog
		err   error
	}
	resultCh := make(chan result)

	for index, breed := range breeds {
		go func(index int, breed Breed) {
			d, err := s.dogs(breed)
			if err != nil {
				resultCh <- result{index: index, err: err}
			}
			resultCh <- result{index: index, dogs: d}
		}(index, breed)
	}

	var results []result
	for i := 0; i < len(breeds); i++ {
		results = append(results, <-resultCh)
	}

	// Sort the results by index
	sort.SliceStable(results, func(a, b int) bool {
		return results[a].index < results[b].index
	})

	var dogs []Dog
	for _, res := range results {
		if res.err == nil {
			dogs = append(dogs, res.dogs...)
		}
	}

	s.logger.WithFields(logrus.Fields{
		"dogs count":     len(dogs),
		"breeds count":   len(breeds),
		"requests count": len(breeds) + 1,
	}).Info("Dog request received")

	return dogs, nil
}
