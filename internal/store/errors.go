package store

import (
	"errors"
	"reflect"
)

var ErrRegistrationNotFound = errors.New("registration not found")

// error type containing information on why a registration already exists
type ErrRegistrationAlreadyExists struct {
	Detail string
}

// satisfying the error interface
func (e ErrRegistrationAlreadyExists) Error() string {
	return "existing registration found: " + e.Detail
}

// implementing Is comparision so we can tell whether or not an error is this
// instance type or not instead of just basing it on the error message itself
func (e ErrRegistrationAlreadyExists) Is(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(e)
}
