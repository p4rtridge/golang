package middleware

import (
	"context"
	"order_service/internal/core"
	"order_service/pkg"
	"order_service/services/auth/controller/api"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func RequireAuth(biz api.AuthUseCase) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		token, err := extractTokenFromCookie(c.Cookies("authAccessToken"))
		if err != nil && token == "" {
			token, err = extractTokenFromHeaderString(c.Get("Authorization"))
			if err != nil {
				return pkg.WriteResponse(c, err)
			}
		}

		sub, tid, err := biz.Verify(ctx, token)
		if err != nil {
			return pkg.WriteResponse(c, err)
		}

		c.Locals(core.KeyRequester, core.NewRequester(sub, tid))

		return c.Next()
	}
}

func extractTokenFromCookie(s string) (string, error) {
	if len(s) == 0 {
		return "", core.ErrUnauthorized.WithError("missing access token")
	}

	return s, nil
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", core.ErrUnauthorized.WithError("missing access token")
	}

	return parts[1], nil
}
