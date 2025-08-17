package validator

import "github.com/go-playground/validator/v10"

// GoValidator wraps the go-playground/validator to provide struct validation.
type GoValidator struct {
	validate *validator.Validate
}

// New creates a new instance of GoValidator.
func New() *GoValidator {
	return &GoValidator{
		validate: validator.New(),
	}
}

// Validate checks the provided struct against validation tags (`validate:"..."`) defined on struct fields.
// Returns an error if validation fails, nil otherwise.
func (v *GoValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}
