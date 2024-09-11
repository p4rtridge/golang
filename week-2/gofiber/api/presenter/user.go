package presenter

import (
	"fmt"

	"github.com/partridge1307/gofiber/entity"
)

type GetUsersPresenter struct{}

func (p *GetUsersPresenter) Present(users *[]entity.User, err error) Response {
	if err != nil {
		fmt.Printf("%T: %v\n", err, err)

		return Response{
			Status:  500,
			Message: "Internal server error",
		}
	}

	return Response{
		Status:  200,
		Message: "OK",
		Data:    users,
	}
}

type GetUserPresenter struct{}

func (p *GetUserPresenter) Present(user *entity.User, err error) Response {
	if err != nil {
		fmt.Printf("%T: %v\n", err, err)

		return Response{
			Status:  500,
			Message: "Internal server error",
		}
	}

	return Response{
		Status:  200,
		Message: "OK",
		Data:    user,
	}
}
