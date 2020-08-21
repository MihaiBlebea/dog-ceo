package cache

import (
	"sync"
	"time"

	"github.com/MihaiBlebea/dog-ceo/dog"
	"github.com/sirupsen/logrus"
)

// Service _
type Service struct {
	dogService dog.Service
	cache      []dog.Dog
	lock       sync.RWMutex
	expiration time.Time
	logger     *logrus.Logger
}

// New returns a new cache service
func New(dog dog.Service, logger *logrus.Logger) *Service {
	s := &Service{
		dogService: dog,
		logger:     logger,
	}

	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		for {
			dogs, err := s.dogService.AllDogs()
			if err != nil {
				continue
			}

			s.lock.Lock()
			s.cache = dogs
			s.extendExpiration()
			s.lock.Unlock()

			s.logger.Info("Worker updated the cache")
			<-ticker.C
		}
	}()

	return s
}

// AllDogs returns all dogs regardless of breed
func (s *Service) AllDogs() ([]dog.Dog, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.cache) > 0 && s.isExpired() == false {
		s.logger.Info("Take dogs from cache")

		return s.cache, nil
	}

	dogs, err := s.dogService.AllDogs()
	if err != nil {
		return []dog.Dog{}, err
	}

	s.cache = dogs
	s.extendExpiration()

	s.logger.Info("Take dogs from api")

	return dogs, nil
}

func (s *Service) isExpired() bool {
	return s.expiration.Before(time.Now())
}

func (s *Service) extendExpiration() {
	s.expiration = time.Now().Add(time.Minute * 5)
}
