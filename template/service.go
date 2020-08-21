package template

// Service interface
type Service interface {
	Load(path string) (*Page, error)
}
