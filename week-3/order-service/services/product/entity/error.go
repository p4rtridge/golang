package entity

import "errors"

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrCannotCreate       = errors.New("cannot create product")
	ErrCannotUpdate       = errors.New("cannot update product")
	ErrCannotDelete       = errors.New("cannot delete product")
	ErrProductNotFound    = errors.New("cannot found product")
)
