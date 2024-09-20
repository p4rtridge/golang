package usecase

import (
	"context"
	"order_service/internal/core"
	orderEntity "order_service/services/order/entity"
	orderRepo "order_service/services/order/repository/postgres"
	productEntity "order_service/services/product/entity"
	userEntity "order_service/services/user/entity"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, data *orderEntity.Order) error
	GetOrders(ctx context.Context) (*[]orderEntity.Order, error)
	GetTopFiveOrdersByPrice(ctx context.Context) (*[]orderEntity.Order, error)
	GetNumOfOrdersByMonth(ctx context.Context) (*[]orderEntity.AggregatedOrdersByMonth, error)
	GetOrder(ctx context.Context, orderId int) (*orderEntity.Order, error)
}

type orderUsecase struct {
	repo orderRepo.OrderRepository
}

func NewUsecase(repo orderRepo.OrderRepository) OrderUsecase {
	return &orderUsecase{
		repo,
	}
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, data *orderEntity.Order) error {
	requester := core.GetRequester(ctx)

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}

	requesterId := int(uid.GetLocalID())
	data.SetUserId(requesterId)

	err = uc.repo.CreateOrder(ctx, data, func(order *orderEntity.Order, user *userEntity.User, products *[]productEntity.Product) (bool, error) {
		// whether any arguments is nil pointer
		if order == nil || user == nil || products == nil {
			return false, orderEntity.ErrInvalidMemory
		}

		// order's items and products must be the same length
		orderItems := order.GetItemsSafe()
		if len(*products) != len(orderItems) {
			return false, orderEntity.ErrNotEqual
		}

		totalPrice := float32(0)

		for idx, item := range orderItems {
			product := &(*products)[idx]

			productQuantity := product.GetQuantity()
			if productQuantity < item.GetQuantity() {
				return false, orderEntity.ErrOutOfStock
			}

			totalPrice += product.GetPrice() * float32(item.GetQuantity())

			// update order's item
			i := (*order).GetItemSafe(idx)
			if i == nil {
				return false, orderEntity.ErrInvalidMemory
			}

			i.SetProductName(product.GetName())
			i.SetProductPrice(product.GetPrice())
			product.SetQuantity(i.GetQuantity())
		}

		if user.GetBalance() < totalPrice {
			return false, orderEntity.ErrInsufficientBalance
		}

		order.SetTotalPrice(totalPrice)
		user.SetBalance(totalPrice)

		return true, nil
	})
	if err != nil {
		if err == orderEntity.ErrOutOfStock {
			return core.ErrConfict.WithError(orderEntity.ErrOutOfStock.Error())
		}
		if err == orderEntity.ErrInsufficientBalance {
			return core.ErrConfict.WithError(orderEntity.ErrInsufficientBalance.Error())
		}

		return core.ErrInternalServerError.WithError(orderEntity.ErrCannotCreateOrder.Error()).WithDebug(err.Error())
	}

	return nil
}

func (uc *orderUsecase) GetOrders(ctx context.Context) (*[]orderEntity.Order, error) {
	requester := core.GetRequester(ctx)

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return nil, core.ErrInternalServerError.WithDebug(err.Error())
	}
	requesterId := int(uid.GetLocalID())

	orders, err := uc.repo.GetOrders(ctx, requesterId)
	if err != nil {
		return nil, core.ErrNotFound.WithError(orderEntity.ErrOrderNotFound.Error()).WithDebug(err.Error())
	}

	return orders, nil
}

func (uc *orderUsecase) GetTopFiveOrdersByPrice(ctx context.Context) (*[]orderEntity.Order, error) {
	orders, err := uc.repo.GetTopFiveOrdersByPrice(ctx)
	if err != nil {
		return nil, core.ErrNotFound.WithError(orderEntity.ErrOrderNotFound.Error()).WithDebug(err.Error())
	}

	return orders, nil
}

func (uc *orderUsecase) GetNumOfOrdersByMonth(ctx context.Context) (*[]orderEntity.AggregatedOrdersByMonth, error) {
	requester := core.GetRequester(ctx)

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return nil, core.ErrInternalServerError.WithDebug(err.Error())
	}
	requesterId := int(uid.GetLocalID())

	orders, err := uc.repo.GetNumOfOrdersPerMonth(ctx, requesterId)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(orderEntity.ErrOrderNotFound.Error()).WithDebug(err.Error())
	}

	return orders, nil
}

func (uc *orderUsecase) GetOrder(ctx context.Context, orderId int) (*orderEntity.Order, error) {
	requester := core.GetRequester(ctx)

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return nil, core.ErrInternalServerError.WithDebug(err.Error())
	}
	requesterId := int(uid.GetLocalID())

	order, err := uc.repo.GetOrder(ctx, requesterId, orderId)
	if err != nil {
		return nil, core.ErrNotFound.WithError(orderEntity.ErrOrderNotFound.Error()).WithDebug(err.Error())
	}

	return order, nil
}
