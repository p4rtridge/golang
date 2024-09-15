package entity

import "errors"

var (
	ErrCannotGetUser    = errors.New("can not get user info")
	ErrCannotAddBalance = errors.New("can not add balance")
)
