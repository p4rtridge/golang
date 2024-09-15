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
	Login(ctx context.Context, data *entity.AuthLogin) (*entity.TokenResponse, error)
	Verify(ctx context.Context, token string) (string, string, error)
	Refresh(ctx context.Context, data *entity.RefreshTokenRequest) (*entity.TokenResponse, error)
	SignOut(ctx context.Context, data *entity.AuthSignOut) error
	SignOutAll(ctx context.Context) error
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

	var data entity.AuthLogin

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	response, err := api.usecase.Login(ctx, &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(response))
}

func (api *api) Refresh(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var data entity.RefreshTokenRequest

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	response, err := api.usecase.Refresh(ctx, &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(response))
}

func (api *api) SignOut(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}
	ctx = core.ContextWithRequester(ctx, requester)

	var data entity.AuthSignOut

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	err := api.usecase.SignOut(ctx, &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(core.ResponseData(true))
}

func (api *api) SignOutAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}
	ctx = core.ContextWithRequester(ctx, requester)

	err := api.usecase.SignOutAll(ctx)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(core.ResponseData(true))
}
