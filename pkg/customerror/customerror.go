package customerror

import (
	"errors"
	"fmt"
)

// ErrorType is the type of an error.
type ErrorType int

type customError struct {
	ErrorType ErrorType
	OrigError error
}

const (
	// NoType error.
	NoType ErrorType = iota
	// BadRequest error.
	BadRequest
	// Unauthorized error.
	Unauthorized
	// NotFound error.
	NotFound
	// Conflict error.
	Conflict
	//UnprocessableEntity error.
	UnprocessableEntity
)

// New is the function that creates a non-type error.
func New(msg string) error {
	return customError{ErrorType: NoType, OrigError: errors.New(msg)}
}

// Newf is the function that creates a non-type error with formatted message.
func Newf(msg string, args ...interface{}) error {
	return customError{ErrorType: NoType, OrigError: fmt.Errorf(msg, args...)}
}

// Error is the function that returns the message of a customError.
func (error customError) Error() string {
	return error.OrigError.Error()
}

// New is the function that creates a new error using an error type.
func (errorType ErrorType) New(msg string) error {
	return customError{ErrorType: errorType, OrigError: errors.New(msg)}
}

// Newf is the function that creates a new error using an error type with formatted message.
func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	return customError{ErrorType: errorType, OrigError: fmt.Errorf(msg, args...)}
}

// GetType is the function that gets the type of the error.
func GetType(err error) ErrorType {
	if customError, ok := err.(customError); ok {
		return customError.ErrorType
	}

	return NoType
}
