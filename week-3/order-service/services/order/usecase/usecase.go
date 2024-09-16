package usecase

import (
	"context"
	"fmt"
	"order_service/internal/core"
	orderEntity "order_service/services/order/entity"
	orderRepo "order_service/services/order/repository/postgres"
	productEntity "order_service/services/product/entity"
	userEntity "order_service/services/user/entity"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, data *orderEntity.OrderRequest) error
	GetOrders(ctx context.Context) (*[]orderEntity.Order, error)
}

type orderUsecase struct {
	repo orderRepo.OrderRepository
}

func NewUsecase(repo orderRepo.OrderRepository) OrderUsecase {
	return &orderUsecase{
		repo,
	}
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, data *orderEntity.OrderRequest) error {
	if err := data.Validate(); err != nil {
		return core.ErrBadRequest.WithError(orderEntity.ErrItemEmpty.Error())
	}

	requester := core.GetRequester(ctx)

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}
	requesterId := int(uid.GetLocalID())

	newItems := make([]orderEntity.OrderItem, 0, len(data.Items))
	for _, reqItem := range data.Items {
		newItem := orderEntity.NewOrderItem(0, reqItem.ProductId, "", 0.0, reqItem.Quantity)

		newItems = append(newItems, newItem)
	}

	newOrder := orderEntity.NewOrder(0, requesterId, 0.0, newItems)

	err = uc.repo.CreateOrder(ctx, &newOrder, func(order *orderEntity.Order, user *userEntity.User, products *[]productEntity.Product) (bool, error) {
		totalPrice := float32(0)

		for idx, item := range order.Items {
			productQuantity := (*products)[idx].Quantity

			if productQuantity < item.Quantity {
				return false, orderEntity.ErrOutOfStock
			}

			totalPrice += (*products)[idx].Price * float32(item.Quantity)
			fmt.Println((*products)[idx].Name)
			order.Items[idx].SetProductName((*products)[idx].Name)
			order.Items[idx].SetProductPrice((*products)[idx].Price)
			(*products)[idx].SetQuantity(productQuantity - item.Quantity)
		}

		if user.Balance < totalPrice {
			return false, orderEntity.ErrInsufficientBalance
		}

		order.SetTotalPrice(totalPrice)
		user.SetBalance(user.Balance - totalPrice)

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
	orders, err := uc.repo.GetOrders(ctx)
	if err != nil {
		return nil, core.ErrNotFound.WithError(orderEntity.ErrOrderNotFound.Error()).WithDebug(err.Error())
	}

	return orders, nil
}
