package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/partridge1307/gofiber/api/middleware"
	"github.com/partridge1307/gofiber/api/presenter"
	"github.com/partridge1307/gofiber/entity"
	"github.com/partridge1307/gofiber/usecase/auth"
)

type AuthHandler struct {
	uc auth.AuthUsecase
}

type AuthRequest struct {
	Username string
	Password string
}

func NewAuthHandler(route fiber.Router, uc auth.AuthUsecase) {
	handler := &AuthHandler{
		uc,
	}

	route.Post("/sign-up", handler.signUp)
	route.Post("/sign-in", handler.signIn)
	route.Get("/private", middleware.JWTMiddleware(handler.uc), func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}

func (h *AuthHandler) signUp(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &AuthRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user := entity.NewUser(0, req.Username, req.Password, 0.0)

	err := h.uc.SignUp(ctx, user)

	presenter := presenter.SignUpPresenter{}
	response := presenter.Present(err)

	return c.Status(response.Status).JSON(fiber.Map{
		"message": response.Message,
		"data":    response.Data,
	})
}

func (h *AuthHandler) signIn(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &AuthRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user := entity.NewUser(0, req.Username, req.Password, 0.0)

	token, err := h.uc.SignIn(ctx, user)

	presenter := presenter.SignInPresenter{}
	response := presenter.Present(token, err)

	c.Cookie(&fiber.Cookie{
		Name:     "auth",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Hour),
		Secure:   false,
		HTTPOnly: true,
	})

	return c.Status(response.Status).JSON(fiber.Map{
		"message": response.Message,
		"data":    response.Data,
	})
}
