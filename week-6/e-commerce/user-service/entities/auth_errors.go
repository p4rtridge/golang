package entities

import "errors"

var (
	ErrCannotAdd       = errors.New("can not add new auth")
	ErrUsernameExisted = errors.New("username has existed")
	ErrCannotRegister  = errors.New("cannot register")
	ErrLoginFailed     = errors.New("username and password is not valid")
	ErrRefreshFailed   = errors.New("cannot refresh token")
	ErrSignoutFailed   = errors.New("sign out failed")
)
