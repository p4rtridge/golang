package pkg

import (
	"order_service/internal/core"

	"github.com/gofiber/fiber/v2"
)

func WriteResponse(c *fiber.Ctx, err error) error {
	if errSt, ok := err.(core.StatusCodeCarrier); ok {
		return c.Status(errSt.StatusCode()).JSON(errSt)
	}

	return c.Status(fiber.StatusInternalServerError).JSON(core.ErrInternalServerError.WithError(err.Error()))
}
