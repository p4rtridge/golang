package core

import (
	"errors"
	"net/http"
)

type DefaultError struct {
	// debug information
	DebugField string `json:"debug,omitempty"`
	// error message
	// ex: the resource could be not found
	ErrorField string `json:"message"`
	// status description
	// ex: not found
	StatusField string `json:"status,omitempty"`
	// status code
	// ex: 404
	CodeField int `json:"code,omitempty"`
}

func (e DefaultError) Error() string {
	return e.ErrorField
}

func (e DefaultError) StatusCode() int {
	return e.CodeField
}

func (e DefaultError) WithError(message string) DefaultError {
	e.ErrorField = message
	return e
}

func (e DefaultError) WithDebug(debug string) DefaultError {
	e.DebugField = debug
	return e
}

type StatusCodeCarrier interface {
	StatusCode() int
}

var ErrInternalServerError = DefaultError{
	StatusField: http.StatusText(http.StatusInternalServerError),
	ErrorField:  "An internal server error occurred, please contact the system administrator",
	CodeField:   http.StatusInternalServerError,
}

var ErrBadRequest = DefaultError{
	StatusField: http.StatusText(http.StatusBadRequest),
	ErrorField:  "The request was malformed or contained invalid parameters",
	CodeField:   http.StatusBadRequest,
}

var ErrUnauthorized = DefaultError{
	StatusField: http.StatusText(http.StatusUnauthorized),
	ErrorField:  "The request could not be authorized",
	CodeField:   http.StatusUnauthorized,
}

var ErrNotFound = DefaultError{
	StatusField: http.StatusText(http.StatusNotFound),
	ErrorField:  "The requested resource could not be found",
	CodeField:   http.StatusNotFound,
}

var ErrConfict = DefaultError{
	StatusField: http.StatusText(http.StatusConflict),
	ErrorField:  "The resource could not be created due to a conflict",
	CodeField:   http.StatusConflict,
}

var ErrRecordNotFound = errors.New("record not found")
