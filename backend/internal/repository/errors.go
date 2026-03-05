// Package repository implements data access layer
package repository

import (
	"errors"
	"fmt"
)

// Common repository errors
var (
	// ErrNotFound indicates the requested resource was not found
	ErrNotFound = errors.New("resource not found")
	// ErrDuplicateKey indicates a unique constraint violation
	ErrDuplicateKey = errors.New("resource already exists")
	// ErrValidation indicates input validation failed
	ErrValidation = errors.New("validation error")
)

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// WrapNotFoundError wraps an error with context
func WrapNotFoundError(err error, resource string) error {
	return fmt.Errorf("%w: %s", ErrNotFound, resource)
}
