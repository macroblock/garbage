package errors

import "fmt"

// TErrors -
type TErrors struct {
	errors []error
}

// NewErrors -
func NewErrors() *TErrors {
	return &TErrors{}
}

// Add -
func (o *TErrors) Add(err ...error) {
	o.errors = append(o.errors, err...)
}

// Addf -
func (o *TErrors) Addf(format string, vals ...interface{}) {
	err := fmt.Errorf(format, vals...)
	o.errors = append(o.errors, err)
}

// Get -
func (o *TErrors) Get() []error {
	if o == nil || len(o.errors) == 0 {
		return nil
	}
	return []error(o.errors)
}
