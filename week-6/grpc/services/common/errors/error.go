package errors

import "google.golang.org/grpc/codes"

type CodeCarrier interface {
	Code() codes.Code
}

type DefaultError struct {
	message string
	code    codes.Code
}

func (e DefaultError) Error() string {
	return e.message
}

func (e DefaultError) Code() codes.Code {
	return e.code
}

var ErrExisted = DefaultError{
	code:    codes.AlreadyExists,
	message: "Already Exists",
}
