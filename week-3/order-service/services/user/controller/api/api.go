package api

import (
	"context"
	"order_service/internal/core"
	"order_service/pkg"
	"order_service/services/user/entity"

	"github.com/gofiber/fiber/v2"
)

type UserUsecase interface {
	GetUsers(ctx context.Context) (*[]entity.User, error)
	GetUser(ctx context.Context, userID int) (*entity.User, error)
	GetUserProfile(ctx context.Context) (*entity.User, error)
	AddUserBalance(ctx context.Context, data *entity.UserRequest) error
}

type api struct {
	usecase UserUsecase
}

func NewAPI(uc UserUsecase) *api {
	return &api{
		usecase: uc,
	}
}

func (api *api) GetUserProfile(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}

	ctx = core.ContextWithRequester(ctx, requester)

	user, err := api.usecase.GetUserProfile(ctx)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(user))
}

func (api *api) GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	users, err := api.usecase.GetUsers(ctx)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(users))
}

func (api *api) GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetId, err := c.ParamsInt("userID")
	if err != nil {
		return pkg.WriteResponse(c, core.ErrNotFound)
	}

	user, err := api.usecase.GetUser(ctx, targetId)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(user))
}

func (api *api) AddUserBalance(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}

	ctx = core.ContextWithRequester(ctx, requester)

	var data entity.UserRequest

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	err := api.usecase.AddUserBalance(ctx, &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(true))
}
