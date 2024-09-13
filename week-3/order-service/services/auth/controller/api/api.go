package api

import (
	"context"
	"order_service/internal/core"
	"order_service/pkg"
	"order_service/services/auth/entity"

	"github.com/gofiber/fiber/v2"
)

type AuthUseCase interface {
	Register(ctx context.Context, data *entity.AuthUsernamePassword) error
	Login(ctx context.Context, data *entity.AuthUsernamePassword) (*entity.TokenResponse, error)
	Verify(ctx context.Context, accessToken string) (string, string, error)
}

type api struct {
	usecase AuthUseCase
}

func NewAPI(uc AuthUseCase) *api {
	return &api{
		usecase: uc,
	}
}

func (api *api) Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var data entity.AuthUsernamePassword

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	err := api.usecase.Register(ctx, &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(core.ResponseData(true))
}

func (api *api) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var data entity.AuthUsernamePassword

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	response, err := api.usecase.Login(ctx, &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(response))
}
