package service

// Service interface
type Service interface {
	Exec(val ...interface{}) error
	Name() string
}
