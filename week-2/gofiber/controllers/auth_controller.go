package controllers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/partridge1307/gofiber/entities"
	"github.com/partridge1307/gofiber/helpers"
	"github.com/partridge1307/gofiber/usecases"
)

type authController struct {
	authUsecase usecases.AuthUsecase
}

type (
	signUpRequest struct {
		Username        string `json:"username" validate:"required"`
		Password        string `json:"password" validate:"required,eqfield=ConfirmPassword"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}
	signInRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

var validate = validator.New()

func NewAuthController(route fiber.Router, authUsecase usecases.AuthUsecase) {
	handler := &authController{
		authUsecase,
	}

	route.Post("/sign-up", handler.signUp)
	route.Post("/sign-in", handler.signIn)
}

func (a *authController) validate(o interface{}, c *fiber.Ctx) error {
	if err := c.BodyParser(o); err != nil {
		return err
	}

	if errs := validate.Struct(o); errs != nil {
		errMsgs := make([]string, 0)
		for _, err := range errs.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("'%s' needs to implement '%s'",
				helpers.GetJsonField(signUpRequest{}, err.StructField()),
				err.Tag(),
			))
		}

		return errors.New(strings.Join(errMsgs, ",\n"))
	}

	return nil
}

func (a *authController) signUp(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &signUpRequest{}

	err := a.validate(req, c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	user := &entities.User{
		Username: req.Username,
		Password: req.Password,
	}

	err = a.authUsecase.SignUp(ctx, user)
	if err != nil {
		if err.Error() == "username already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "fail",
				"message": "user already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "ok",
		"data":   nil,
	})
}

func (a *authController) signIn(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &signInRequest{}

	err := a.validate(req, c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	user := &entities.User{
		Username: req.Username,
		Password: req.Password,
	}

	token, err := a.authUsecase.SignIn(ctx, user)
	if err != nil {
		if err.Error() == "unmatched" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "username or password does not match",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"token":  token,
	})
}
