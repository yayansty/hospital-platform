package helpers

type ValidationErrors map[string]string

type ValidationError struct {
	Errors ValidationErrors
}

func (e *ValidationError) Error() string {
	return "validation failed"
}
