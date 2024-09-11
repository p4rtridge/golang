package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/partridge1307/gofiber/api/presenter"
	"github.com/partridge1307/gofiber/usecase/user"
)

type UserHandler struct {
	uc user.UserUsecase
}

func NewUserHandler(route fiber.Router, uc user.UserUsecase) {
	handler := &UserHandler{
		uc,
	}

	route.Get("/", handler.getUsers)
	route.Get("/:userID", handler.getUser)
}

func (h *UserHandler) getUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	users, err := h.uc.GetUsers(ctx)

	presenter := presenter.GetUsersPresenter{}
	response := presenter.Present(users, err)

	return c.Status(response.Status).JSON(fiber.Map{
		"message": response.Message,
		"data":    response.Data,
	})
}

func (h *UserHandler) getUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetUserID, err := c.ParamsInt("userID")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	user, err := h.uc.GetUser(ctx, targetUserID)

	presenter := presenter.GetUserPresenter{}
	response := presenter.Present(user, err)

	return c.Status(response.Status).JSON(fiber.Map{
		"message": response.Message,
		"data":    response.Data,
	})
}
