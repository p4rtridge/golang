package entity

import "errors"

var (
	ErrMissingField        = errors.New("missing item's field")
	ErrInvalidMemory       = errors.New("invalid memory in required variable")
	ErrNotEqual            = errors.New("products and order's items is not equal")
	ErrItemEmpty           = errors.New("item cannot be empty")
	ErrCannotCreateOrder   = errors.New("order cannot be create")
	ErrInsufficientBalance = errors.New("order cannot be create because user's balance is insufficient")
	ErrOutOfStock          = errors.New("one item in order's items is out of stock")
	ErrOrderNotFound       = errors.New("cannot be found any orders")
)
