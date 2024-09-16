package entity

import "errors"

var (
	ErrItemEmpty         = errors.New("item cannot be empty")
	ErrCannotCreateOrder = errors.New("order cannot be create")
)
