package order

import "errors"

var (
	ErrInvalidJSON = errors.New("invalid JSON")
	ErrNilOrder    = errors.New("nil order")
	ErrValidation  = errors.New("failed to validate order")
	ErrCreateOrder = errors.New("failed to create order")
)
