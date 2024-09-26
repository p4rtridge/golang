package api

import (
	"order_service/internal/core"
	"order_service/pkg"
	"order_service/services/auth/entity"
	authUc "order_service/services/auth/usecase"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Register(*fiber.Ctx) error
	Login(*fiber.Ctx) error
	Refresh(*fiber.Ctx) error
	SignOut(*fiber.Ctx) error
	SignOutAll(*fiber.Ctx) error
}

type service struct {
	usecase authUc.AuthUseCase
}

func NewService(uc authUc.AuthUseCase) AuthService {
	return &service{
		usecase: uc,
	}
}

func (srv *service) Register(c *fiber.Ctx) error {
	var data entity.AuthUsernamePassword

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	if err := data.Validate(); err != nil {
		return pkg.WriteResponse(c, err)
	}

	err := srv.usecase.Register(c.Context(), data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(core.ResponseData(true))
}

func (srv *service) Login(c *fiber.Ctx) error {
	var data entity.AuthLogin

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	if err := data.Validate(); err != nil {
		return pkg.WriteResponse(c, err)
	}

	response, err := srv.usecase.Login(c.Context(), data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(response))
}

func (srv *service) Refresh(c *fiber.Ctx) error {
	var data entity.RefreshTokenRequest

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	response, err := srv.usecase.Refresh(c.Context(), data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(response))
}

func (srv *service) SignOut(c *fiber.Ctx) error {
	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}
	ctx := core.ContextWithRequester(c.Context(), requester)

	var data entity.AuthSignOut

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	err := srv.usecase.SignOut(ctx, data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(core.ResponseData(true))
}

func (srv *service) SignOutAll(c *fiber.Ctx) error {
	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}
	ctx := core.ContextWithRequester(c.Context(), requester)

	err := srv.usecase.SignOutAll(ctx)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(core.ResponseData(true))
}
