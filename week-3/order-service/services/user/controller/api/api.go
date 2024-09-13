package api

import (
	"context"
	"order_service/internal/core"
	"order_service/pkg"
	"order_service/services/user/entity"

	"github.com/gofiber/fiber/v2"
)

type UserUsecase interface {
	GetUserProfile(ctx context.Context) (*entity.User, error)
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
