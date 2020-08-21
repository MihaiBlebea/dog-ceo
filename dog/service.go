package dog

// Service interface
type Service interface {
	AllDogs() ([]Dog, error)
}
