package api

import (
	"order_service/internal/core"
	"order_service/pkg"
	"order_service/services/user/entity"
	userUc "order_service/services/user/usecase"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	GetUsers(*fiber.Ctx) error
	GetUser(*fiber.Ctx) error
	GetUserProfile(*fiber.Ctx) error
	AddUserBalance(*fiber.Ctx) error
}

type service struct {
	usecase userUc.UserUsecase
}

func NewService(uc userUc.UserUsecase) UserService {
	return &service{
		usecase: uc,
	}
}

// Get User Profile godoc
// @summary Get User Profile
// @description Get the current user profile
// @tags users
// @security BearerAuth
// @success 200 {object} entity.User
// @failure 401 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /users/profile [get]
func (srv *service) GetUserProfile(c *fiber.Ctx) error {
	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return pkg.WriteResponse(c, core.ErrInternalServerError.WithDebug(err.Error()))
	}
	requesterId := int(uid.GetLocalID())

	user, err := srv.usecase.GetUserProfile(c.Context(), requesterId)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(user))
}

// Get Users godoc
// @summary Get Users
// @description Get entire users
// @tags users
// @security BearerAuth
// @success 200 {array} entity.User
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /users/ [get]
func (srv *service) GetUsers(c *fiber.Ctx) error {
	users, err := srv.usecase.GetUsers(c.Context())
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(users))
}

// Get User godoc
// @summary Get User
// @description Get specific user
// @tags users
// @security BearerAuth
// @param userID path int true "User's ID"
// @success 200 {object} entity.User
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /users/:userID [get]
func (srv *service) GetUser(c *fiber.Ctx) error {
	targetId, err := c.ParamsInt("userID")
	if err != nil {
		return pkg.WriteResponse(c, core.ErrNotFound)
	}

	user, err := srv.usecase.GetUser(c.Context(), targetId)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(user))
}

// Add User Balance godoc
// @summary Add User Balance
// @description Add user balance
// @tags users
// @security BearerAuth
// @param payload body entity.UserRequest true "User request body"
// @success 200
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /users/balance [post]
func (srv *service) AddUserBalance(c *fiber.Ctx) error {
	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return pkg.WriteResponse(c, core.ErrInternalServerError.WithDebug(err.Error()))
	}

	requesterId := int(uid.GetLocalID())
	var data entity.UserRequest

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	err = srv.usecase.AddUserBalance(c.Context(), requesterId, data.Balance)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(true))
}
