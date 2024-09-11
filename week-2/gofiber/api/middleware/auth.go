package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/partridge1307/gofiber/usecase/auth"
)

func JWTMiddleware(authUsecase auth.AuthUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		authToken := c.Cookies("auth")
		if authToken == "" {
			return c.Status(401).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		ok, err := authUsecase.Verify(ctx, authToken)
		if err != nil && !ok {
			return c.Status(401).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		return c.Next()
	}
}
