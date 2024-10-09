package api

import (
	"fmt"
	"order_service/internal/core"
	"order_service/pkg"
	orderEntity "order_service/services/order/entity"
	orderUsecase "order_service/services/order/usecase"

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

// Create Order godoc
// @summary Create a new order
// @description Create a new order with the input payload
// @tags orders
// @accept application/json
// @security BearerAuth
// @param payload body entity.OrderRequest true "Create order request body"
// @success 201
// @failure 400 {object} core.DefaultError
// @failure 401 {object} core.DefaultError
// @failure 409 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /orders/ [post]
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

// Get Orders godoc
// @summary Get Orders
// @description Get all of orders of the current user or entire users based on user's role
// @tags orders
// @security BearerAuth
// @success 200 {array} entity.Order
// @failure 401 {object} core.DefaultError
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /orders/ [get]
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

// Get Orders Summarize godoc
// @summary Get Orders Summarize
// @description Get summarized orders of the current user and export to excel
// @tags orders
// @security BearerAuth
// @param payload body entity.OrdersSummarizeReq true "Order summary request body"
// @success 301
// @failure 500 {object} core.DefaultError
// @router /orders/summarize [get]
func (srv *service) GetOrdersSummarize(c *fiber.Ctx) error {
	var orderSummaryReq orderEntity.OrdersSummarizeReq

	if err := c.BodyParser(&orderSummaryReq); err != nil {
		return pkg.WriteResponse(c, err)
	}

	datas, err := srv.usecase.GetOrdersSummarize(c.Context(), orderSummaryReq.StartDate, orderSummaryReq.EndDate)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	excelName, err := pkg.GenerateExcel(datas, orderSummaryReq.StartDate, orderSummaryReq.EndDate)
	if err != nil {
		return pkg.WriteResponse(c, core.ErrInternalServerError.WithDebug(err.Error()))
	}

	return c.Redirect(fmt.Sprintf("/static/%s", excelName), fiber.StatusMovedPermanently)
}

// Get Top Five Orders Order By Price godoc
// @summary Get Top Five Orders
// @description Get top five orders order by price
// @tags orders
// @security BearerAuth
// @success 200 {array} entity.Order
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /orders/top-by-price [get]
func (srv *service) GetTopFiveOrdersByPrice(c *fiber.Ctx) error {
	orders, err := srv.usecase.GetTopFiveOrdersByPrice(c.Context())
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(orders))
}

// Get Order godoc
// @summary Get Order
// @description Get specific order of the current user and export to pdf
// @tags orders
// @security BearerAuth
// @param orderID path int true "Order's ID"
// @success 301
// @failure 401 {object} core.DefaultError
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /orders/:orderID/invoice [get]
func (srv *service) GetOrder(c *fiber.Ctx) error {
	targetOrderId, err := c.ParamsInt("orderID")
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return pkg.WriteResponse(c, core.ErrInternalServerError.WithDebug(err.Error()))
	}
	requesterId := int(uid.GetLocalID())

	order, err := srv.usecase.GetOrder(c.Context(), requesterId, targetOrderId)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	err = pkg.GeneratePDF(order)
	if err != nil {
		return pkg.WriteResponse(c, core.ErrInternalServerError.WithDebug(err.Error()))
	}

	return c.Redirect(fmt.Sprintf("/static/invoice-%d-%d.pdf", order.GetUserIdSafe(), order.GetIdSafe()), fiber.StatusMovedPermanently)
}

// Get Aggregated Orders By Month godoc
// @summary Get Aggregated Orders By Month
// @description Get aggregated orders group by month of the current user
// @tags orders
// @security BearerAuth
// @success 200
// @failure 401 {object} core.DefaultError
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /orders/orders-by-month [get]
func (srv *service) GetNumOfOrdersByMonth(c *fiber.Ctx) error {
	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return pkg.WriteResponse(c, core.ErrInternalServerError.WithDebug(err.Error()))
	}
	requesterId := int(uid.GetLocalID())

	orders, err := srv.usecase.GetNumOfOrdersByMonth(c.Context(), requesterId)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(orders))
}
