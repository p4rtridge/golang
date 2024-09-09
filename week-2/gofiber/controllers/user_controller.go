package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/partridge1307/gofiber/usecases"
)

type userController struct {
	uc usecases.UserUsecase
}

func NewUserController(route fiber.Router, uc usecases.UserUsecase) {
	handler := &userController{
		uc,
	}

	route.Get("/", handler.getUsers)
	route.Get("/:userID", handler.getUser)
}

func (u *userController) getUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	users, err := u.uc.GetUsers(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}

func (u *userController) getUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetUserID, err := c.ParamsInt("userID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Please specify a valid user id",
		})
	}

	user, err := u.uc.GetUser(ctx, targetUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}
