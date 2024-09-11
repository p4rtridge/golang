package auth

import "strings"

type UserExistError struct{}

func (e *UserExistError) Error() string {
	return "user already exists"
}

type UserUnmatchError struct{}

func (e *UserUnmatchError) Error() string {
	return "user unmatch"
}

type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Errors, ",\n")
}
