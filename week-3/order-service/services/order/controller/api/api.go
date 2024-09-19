package api

import (
	"order_service/internal/core"
	"order_service/pkg"
	orderEntity "order_service/services/order/entity"
	orderUsecase "order_service/services/order/usecase"

	"github.com/gofiber/fiber/v2"
)

type OrderService interface {
	CreateOrder(*fiber.Ctx) error
	GetOrders(*fiber.Ctx) error
}

type service struct {
	usecase orderUsecase.OrderUsecase
}

func NewService(uc orderUsecase.OrderUsecase) OrderService {
	return &service{
		usecase: uc,
	}
}

func (srv *service) CreateOrder(c *fiber.Ctx) error {
	var data orderEntity.OrderRequest

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	if err := data.Validate(); err != nil {
		return core.ErrBadRequest.WithError(orderEntity.ErrItemEmpty.Error())
	}

	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}

	ctx := core.ContextWithRequester(c.Context(), requester)

	dataItems := data.GetItems()
	newItems := make([]orderEntity.OrderItem, 0, len(dataItems))

	for _, reqItem := range dataItems {
		newItem := orderEntity.NewOrderItem(0, reqItem.GetItemId(), "", 0.0, reqItem.GetItemQuantity())

		newItems = append(newItems, newItem)
	}

	newOrder := orderEntity.NewOrder(0, 0, 0.0, newItems)

	err := srv.usecase.CreateOrder(ctx, &newOrder)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(core.ResponseData(true))
}

func (srv *service) GetOrders(c *fiber.Ctx) error {
	orders, err := srv.usecase.GetOrders(c.Context())
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(orders))
}
