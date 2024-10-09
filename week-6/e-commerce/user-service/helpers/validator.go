package helpers

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

var (
	ErrUsernameIsNotValid = errors.New("username is not valid. Accept only alphanumeric characters and hyphen letter")
	ErrPasswordIsTooShort = errors.New("password too short, 8 characters at least")
	ErrPasswordIsNotValid = errors.New("password is not valid. Must start with an uppercase or lowercase letter and contain at least one uppercase, lowercase and one digit")
	ErrDeviceIdIsNotValid = errors.New("device id is not valid. Must be uuid")
)

// Check whether username is valid
//
// username must starts with a lowercase chatacter
func UsernameIsValid(username string) error {
	re := regexp.MustCompile("^[a-z][-a-z0-9]+$")

	if !re.MatchString(username) {
		return ErrUsernameIsNotValid
	}

	return nil
}

// Check whether password is valid
//
// password must at least 8 characters, one lowercase letter, one uppercase letter, one digit
func CheckPassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordIsTooShort
	}

	// ^[A-Za-z0-9]* : Ensure the password starts with any alphanumeric letter
	// [A-Z][A-Za-z0-9]* : Ensure at least one uppercase letter
	// [a-z][A-Za-z0-9]* : Ensure at least one lowercase letter
	// [0-9][A-Za-z0-9]* : Ensure at least one digit

	// comment these line to load test
	re := regexp.MustCompile("^[A-Za-z0-9]*[A-Z][A-Za-z0-9]*[a-z][A-Za-z0-9]*[0-9][A-Za-z0-9]*")

	if !re.MatchString(password) {
		return ErrPasswordIsNotValid
	}

	return nil
}

// Check whether device id is uuid
func DeviceIdIsValid(u string) error {
	_, err := uuid.Parse(u)
	if err != nil {
		return ErrDeviceIdIsNotValid
	}

	return nil
}
