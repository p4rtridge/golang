package handler

import (
	"context"
	"kitchen/services/common/errors"
	"kitchen/services/common/genproto/orders"
	"kitchen/services/orders/entity"
	"kitchen/services/orders/service"
)

type OrdersRPCHandler struct {
	orders.UnimplementedOrderServiceServer
	service service.OrderService
}

func NewOrdersRPCHandler(srv service.OrderService) orders.OrderServiceServer {
	rpcHandler := &OrdersRPCHandler{
		service: srv,
	}

	return rpcHandler
}

func (h *OrdersRPCHandler) CreateOrder(
	ctx context.Context,
	in *orders.CreateOrderRequest,
) (*orders.CreateOrderResponse, error) {
	factory := entity.OrderFactory{}

	err := h.service.CreateOrder(ctx, factory.CreateFromProto(in))
	if err != nil {
		return nil, errors.WriteError(err)
	}

	return &orders.CreateOrderResponse{
		Status: "success",
	}, nil
}

func (h *OrdersRPCHandler) GetOrders(
	ctx context.Context,
	in *orders.GetOrdersRequest,
) (*orders.GetOrdersResponse, error) {
	o, err := h.service.GetOrders(ctx, 1)
	if err != nil {
		return nil, errors.WriteError(err)
	}

	orderResp := make([]*orders.Order, 0)
	for _, or := range o {
		orderResp = append(
			orderResp,
			&orders.Order{
				OrderID:    or.OrderID,
				CustomerID: or.CustomerID,
				ProductID:  or.ProductID,
				Quantity:   or.Quantity,
			},
		)
	}

	return &orders.GetOrdersResponse{
		Orders: orderResp,
	}, nil
}
