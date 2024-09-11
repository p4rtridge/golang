package presenter

import (
	"errors"
	"fmt"

	"github.com/partridge1307/gofiber/usecase/auth"
)

type SignUpPresenter struct{}

func (p *SignUpPresenter) Present(err error) Response {
	if err != nil {
		fmt.Printf("%T: %v\n", err, err)

		var validationErr *auth.ValidationError
		if errors.As(err, &validationErr) {
			return Response{
				Status:  400,
				Message: validationErr.Error(),
			}
		}

		return Response{
			Status:  500,
			Message: "Internal server error",
		}
	}

	return Response{
		Status:  201,
		Message: "Created",
		Data:    nil,
	}
}

type SignInPresenter struct{}

func (p *SignInPresenter) Present(token string, err error) Response {
	if err != nil {
		fmt.Printf("%T: %v\n", err, err)

		var validationErr *auth.ValidationError
		if errors.As(err, &validationErr) {
			return Response{
				Status:  400,
				Message: validationErr.Error(),
			}
		}

		return Response{
			Status:  500,
			Message: "Internal server error",
		}
	}

	return Response{
		Status:  200,
		Message: "Authorized",
		Data:    token,
	}
}
