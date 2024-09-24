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

func (srv *service) GetUsers(c *fiber.Ctx) error {
	users, err := srv.usecase.GetUsers(c.Context())
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(users))
}

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
