package cache

import (
	"sync"

	"github.com/MihaiBlebea/dog-ceo/dog"
	"github.com/sirupsen/logrus"
)

// Service _
type Service struct {
	DogService dog.Service
	Dogs       []dog.Dog
	lock       sync.RWMutex
	Logger     *logrus.Logger
}

// AllDogs returns all dogs regardless of breed
func (s *Service) AllDogs() ([]dog.Dog, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.Dogs) > 0 {
		s.Logger.Info("Requested all dogs from cache")

		return s.Dogs, nil
	}

	dogs, err := s.DogService.AllDogs()
	if err != nil {
		return []dog.Dog{}, err
	}

	s.Dogs = dogs
	s.Logger.Info("Requested all dogs from api")

	return dogs, nil
}
