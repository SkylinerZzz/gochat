package service

// Service interface
type Service interface {
	Exec(...interface{}) error
	Name() string
}
