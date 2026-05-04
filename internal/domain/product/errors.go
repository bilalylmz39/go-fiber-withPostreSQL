package product

import "errors"

var (
	ErrNotFound   = errors.New("product not found")
	ErrValidation = errors.New("product validation failed")
)
