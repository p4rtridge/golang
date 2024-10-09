package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/partridge1307/service-ctx/core"
)

func WriteResponse(c *fiber.Ctx, err error) error {
	if errSt, ok := err.(core.StatusCodeCarrier); ok {
		return c.Status(errSt.StatusCode()).JSON(errSt)
	}

	return c.Status(fiber.StatusInternalServerError).JSON(core.ErrInternalServerError.WithError(err.Error()))
}
