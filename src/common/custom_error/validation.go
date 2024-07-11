package custom_error

import "fmt"

type ValidationError struct {
	Name    string
	Message string
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Name, err.Message)
}
