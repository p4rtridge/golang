package entity

import "errors"

var (
	ErrCannotCreate = errors.New("cannot create product")
	ErrCannotUpdate = errors.New("cannot update product")
	ErrCannotDelete = errors.New("cannot delete product")
)