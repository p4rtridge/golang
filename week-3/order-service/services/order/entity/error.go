package entity

import "errors"

var (
	ErrItemEmpty           = errors.New("item cannot be empty")
	ErrCannotCreateOrder   = errors.New("order cannot be create")
	ErrInsufficientBalance = errors.New("order cannot be create because user's balance is insufficient")
	ErrOutOfStock          = errors.New("one item in order's items is out of stock")
	ErrOrderNotFound       = errors.New("cannot be found any orders")
)
