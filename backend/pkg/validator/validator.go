package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator instance
func New() (*Validator, error) {
	return &Validator{
		validate: validator.New(),
	}, nil
}

// Validate validates the given struct
func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

// Var validates a single variable
func (v *Validator) Var(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}
