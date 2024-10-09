package usecase

import (
	"context"
	"order_service/internal/core"
	orderEntity "order_service/services/order/entity"
	orderRepo "order_service/services/order/repository/postgres"
	productEntity "order_service/services/product/entity"
	userEntity "order_service/services/user/entity"
	"time"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, data *orderEntity.Order) error
	CreateOrderCallback(order *orderEntity.Order, user *userEntity.User, products *[]productEntity.Product) (bool, error)
	GetOrders(ctx context.Context) (*[]orderEntity.Order, error)
	GetTopFiveOrdersByPrice(ctx context.Context) (*[]orderEntity.Order, error)
	GetNumOfOrdersByMonth(ctx context.Context, userId int) (*[]orderEntity.AggregatedOrdersByMonth, error)
	GetOrdersSummarize(ctx context.Context, startDate, endDate time.Time) (*[]orderEntity.OrdersSummarize, error)
	GetOrder(ctx context.Context, userId, orderId int) (*orderEntity.Order, error)
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

	role := uid.GetRole()
	if role == 1 {
		return core.ErrBadRequest.WithError(orderEntity.ErrCannotCreateOrder.Error())
	}

	requesterId := uid.GetLocalID()
	data.SetUserId(int(requesterId))

	err = uc.repo.CreateOrder(ctx, data, uc.CreateOrderCallback)
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

func (uc *orderUsecase) CreateOrderCallback(order *orderEntity.Order, user *userEntity.User, products *[]productEntity.Product) (bool, error) {
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
}

func (uc *orderUsecase) GetOrders(ctx context.Context) (*[]orderEntity.Order, error) {
	requester := core.GetRequester(ctx)
	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return nil, core.ErrInternalServerError.WithDebug(err.Error())
	}

	requesterId := uid.GetLocalID()
	role := uid.GetRole()

	var orders *[]orderEntity.Order

	if role == 1 {
		orders, err = uc.repo.GetOrders(ctx)
	} else {
		orders, err = uc.repo.GetOrdersByUserId(ctx, int(requesterId))
	}

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

func (uc *orderUsecase) GetNumOfOrdersByMonth(ctx context.Context, userId int) (*[]orderEntity.AggregatedOrdersByMonth, error) {
	orders, err := uc.repo.GetNumOfOrdersPerMonth(ctx, userId)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(orderEntity.ErrOrderNotFound.Error()).WithDebug(err.Error())
	}

	return orders, nil
}

func (uc *orderUsecase) GetOrdersSummarize(ctx context.Context, startDate, endDate time.Time) (*[]orderEntity.OrdersSummarize, error) {
	datas, err := uc.repo.GetOrdersSummarize(ctx, startDate, endDate)
	if err != nil {
		return nil, core.ErrInternalServerError.WithDebug(err.Error())
	}

	return datas, nil
}

func (uc *orderUsecase) GetOrder(ctx context.Context, userId, orderId int) (*orderEntity.Order, error) {
	order, err := uc.repo.GetOrder(ctx, userId, orderId)
	if err != nil {
		return nil, core.ErrNotFound.WithError(orderEntity.ErrOrderNotFound.Error()).WithDebug(err.Error())
	}

	return order, nil
}
