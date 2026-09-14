package helpers

import "errors"

type ValidationErrors map[string]string

type ValidationError struct {
	Errors ValidationErrors
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
	ErrDatabase = errors.New("database error")
)
