package api

import (
	"user-service/entities"
	"user-service/helpers"
	"user-service/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/partridge1307/service-ctx/core"
)

type AuthService interface {
	Register(*fiber.Ctx) error
	Login(*fiber.Ctx) error
	Refresh(*fiber.Ctx) error
	SignOut(*fiber.Ctx) error
	SignOutAll(*fiber.Ctx) error
}

type service struct {
	usecase usecases.AuthUsecase
}

func NewService(usecase usecases.AuthUsecase) AuthService {
	return &service{
		usecase,
	}
}

// Register godoc
// @summary Create a new user
// @description Create a new user with the input payload
// @tags auth
// @accept application/json
// @param payload body entities.AuthUsernamePassword true "Auth register request body"
// @success 201
// @failure 400 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /auth/register [post]
func (srv *service) Register(c *fiber.Ctx) error {
	var data entities.AuthUsernamePassword

	if err := c.BodyParser(&data); err != nil {
		return helpers.WriteResponse(c, err)
	}

	if err := data.Validate(); err != nil {
		return helpers.WriteResponse(c, err)
	}

	err := srv.usecase.Register(c.Context(), data)
	if err != nil {
		return helpers.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(core.ResponseData(true))
}

// Login godoc
// @summary Login
// @description Login a specific user with the input payload
// @tags auth
// @accept application/json
// @param payload body entities.AuthLogin true "Auth login request body"
// @success 200 {object} entities.TokenResponse
// @failure 400 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /auth/login [post]
func (srv *service) Login(c *fiber.Ctx) error {
	var data entities.AuthLogin

	if err := c.BodyParser(&data); err != nil {
		return helpers.WriteResponse(c, err)
	}

	if err := data.Validate(); err != nil {
		return helpers.WriteResponse(c, err)
	}

	response, err := srv.usecase.Login(c.Context(), data)
	if err != nil {
		return helpers.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(response))
}

// Refresh godoc
// @summary Refresh token
// @description Rotate current user's session token
// @tags auth
// @accept application/json
// @param payload body entities.RefreshTokenRequest true "Auth refresh request body"
// @success 200 {object} entities.TokenResponse
// @failure 400 {object} core.DefaultError
// @failure 401 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /auth/refresh [post]
func (srv *service) Refresh(c *fiber.Ctx) error {
	var data entities.RefreshTokenRequest

	if err := c.BodyParser(&data); err != nil {
		return helpers.WriteResponse(c, err)
	}

	response, err := srv.usecase.Refresh(c.Context(), data)
	if err != nil {
		return helpers.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(response))
}

// Sign Out godoc
// @summary Sign out
// @description Sign out current user's session with specific device
// @tags auth
// @accept application/json
// @security BearerAuth
// @param payload body entities.AuthSignOut true "Auth sign out request body"
// @success 204
// @failure 401 {object} core.DefaultError
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /auth/sign-out [post]
func (srv *service) SignOut(c *fiber.Ctx) error {
	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return helpers.WriteResponse(c, core.ErrUnauthorized)
	}
	ctx := core.ContextWithRequester(c.Context(), requester)

	var data entities.AuthSignOut

	if err := c.BodyParser(&data); err != nil {
		return helpers.WriteResponse(c, err)
	}

	err := srv.usecase.SignOut(ctx, data)
	if err != nil {
		return helpers.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(core.ResponseData(true))
}

// Sign Out All godoc
// @summary Sign out all
// @description Sign out current user's session from all devices
// @tags auth
// @accept application/json
// @security BearerAuth
// @success 204
// @failure 401 {object} core.DefaultError
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /auth/sign-out-all [post]
func (srv *service) SignOutAll(c *fiber.Ctx) error {
	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return helpers.WriteResponse(c, core.ErrUnauthorized)
	}
	ctx := core.ContextWithRequester(c.Context(), requester)

	err := srv.usecase.SignOutAll(ctx)
	if err != nil {
		return helpers.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(core.ResponseData(true))
}
