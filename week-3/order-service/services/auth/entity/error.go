package entity

import "errors"

var (
	ErrUsernameExisted = errors.New("username has existed")
	ErrCannotRegister  = errors.New("cannot register")
	ErrLoginFailed     = errors.New("username and password is not valid")
)
