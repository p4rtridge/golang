package service

import (
	"context"
	"kitchen/services/common/errors"
	"kitchen/services/orders/entity"
)

type OrderService interface {
	CreateOrder(context.Context, entity.Order) error
	GetOrders(context.Context, int32) ([]*entity.Order, error)
}

var store = make([]*entity.Order, 0)

type orderService struct {
	store []*entity.Order
}

func NewOrderService() OrderService {
	return &orderService{
		store: store,
	}
}

func (srv *orderService) CreateOrder(ctx context.Context, order entity.Order) error {
	for _, o := range srv.store {
		if *o == order {
			return errors.ErrExisted
		}
	}

	order.OrderID += 1
	srv.store = append(srv.store, &order)

	return nil
}

func (srv *orderService) GetOrders(ctx context.Context, customerId int32) ([]*entity.Order, error) {
	return srv.store, nil
}
