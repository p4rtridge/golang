package api

import (
	"fmt"
	"order_service/internal/core"
	"order_service/pkg"
	orderEntity "order_service/services/order/entity"
	orderUsecase "order_service/services/order/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
)

type OrderService interface {
	CreateOrder(*fiber.Ctx) error
	GetOrders(*fiber.Ctx) error
	GetOrdersSummarize(*fiber.Ctx) error
	GetTopFiveOrdersByPrice(*fiber.Ctx) error
	GetNumOfOrdersByMonth(*fiber.Ctx) error
	GetOrder(*fiber.Ctx) error
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

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}

	requesterId := int(uid.GetLocalID())
	dataItems := data.GetItems()
	newItems := make([]orderEntity.OrderItem, 0, len(dataItems))

	for _, reqItem := range dataItems {
		newItem := orderEntity.NewOrderItem(0, reqItem.GetItemId(), "", 0.0, reqItem.GetItemQuantity())

		newItems = append(newItems, newItem)
	}

	newOrder := orderEntity.NewOrder(0, requesterId, 0.0, newItems)

	err = srv.usecase.CreateOrder(c.Context(), &newOrder)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(core.ResponseData(true))
}

func (srv *service) GetOrders(c *fiber.Ctx) error {
	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}

	ctx := core.ContextWithRequester(c.Context(), requester)

	orders, err := srv.usecase.GetOrders(ctx)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(orders))
}

func (srv *service) GetOrdersSummarize(c *fiber.Ctx) error {
	var inlineStruct struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	if err := c.BodyParser(&inlineStruct); err != nil {
		return pkg.WriteResponse(c, err)
	}

	datas, err := srv.usecase.GetOrdersSummarize(c.Context(), inlineStruct.StartDate, inlineStruct.EndDate)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	excelName, err := pkg.GenerateExcel(datas, inlineStruct.StartDate, inlineStruct.EndDate)
	if err != nil {
		return pkg.WriteResponse(c, core.ErrInternalServerError.WithDebug(err.Error()))
	}

	return c.Redirect(fmt.Sprintf("/static/%s", excelName), fiber.StatusMovedPermanently)
}

func (srv *service) GetTopFiveOrdersByPrice(c *fiber.Ctx) error {
	orders, err := srv.usecase.GetTopFiveOrdersByPrice(c.Context())
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(orders))
}

func (srv *service) GetOrder(c *fiber.Ctx) error {
	targetOrderId, err := c.ParamsInt("orderID")
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}

	ctx := core.ContextWithRequester(c.Context(), requester)

	order, err := srv.usecase.GetOrder(ctx, targetOrderId)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	err = pkg.GeneratePDF(order)
	if err != nil {
		return pkg.WriteResponse(c, core.ErrInternalServerError.WithDebug(err.Error()))
	}

	return c.Redirect(fmt.Sprintf("/static/invoice-%d-%d.pdf", order.GetUserIdSafe(), order.GetIdSafe()), fiber.StatusMovedPermanently)
}

func (srv *service) GetNumOfOrdersByMonth(c *fiber.Ctx) error {
	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}

	ctx := core.ContextWithRequester(c.Context(), requester)

	orders, err := srv.usecase.GetNumOfOrdersByMonth(ctx)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(orders))
}
