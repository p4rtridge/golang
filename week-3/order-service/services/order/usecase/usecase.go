package usecase

import (
	"context"
	"order_service/internal/core"
	authEntity "order_service/services/auth/entity"
	orderEntity "order_service/services/order/entity"
	orderRepo "order_service/services/order/repository/postgres"
	productEntity "order_service/services/product/entity"
	productRepo "order_service/services/product/repository/postgres"
	userRepo "order_service/services/user/repository/postgres"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, data *orderEntity.OrderRequest) error
	GetOrders(ctx context.Context) (*[]orderEntity.Order, error)
}

type orderUsecase struct {
	repo        orderRepo.OrderRepository
	userRepo    userRepo.UserRepository
	productRepo productRepo.ProductRepository
}

func NewUsecase(repo orderRepo.OrderRepository, userRepo userRepo.UserRepository, productRepo productRepo.ProductRepository) OrderUsecase {
	return &orderUsecase{
		repo,
		userRepo,
		productRepo,
	}
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, data *orderEntity.OrderRequest) error {
	if err := data.Validate(); err != nil {
		return core.ErrBadRequest.WithError(err.Error())
	}

	requester := core.GetRequester(ctx)

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}
	requesterId := int(uid.GetLocalID())

	user, err := uc.userRepo.GetUserById(ctx, requesterId)
	if err != nil {
		return core.ErrUnauthorized.WithError(authEntity.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	orderItems := make([]orderEntity.OrderItem, 0, len(data.Items))
	products := make([]productEntity.Product, 0, len(data.Items))
	for _, item := range data.Items {
		product, err := uc.productRepo.GetProduct(ctx, item.ProductId)
		if err != nil {
			return core.ErrNotFound.WithError(productEntity.ErrProductNotFound.Error()).WithDebug(err.Error())
		}

		if product.Quantity < item.Quantity {
			return core.ErrConfict.WithError(orderEntity.ErrOutOfStock.Error())
		}

		product.AddQuantity(-item.Quantity)

		orderItems = append(orderItems, orderEntity.NewOrderItem(0, product.Id, product.Name, product.Price, item.Quantity))
		products = append(products, *product)
	}

	newOrder := orderEntity.NewOrder(0, requesterId, 0.0, orderItems)

	orderTotalPrice := newOrder.CalculatePrice()

	if user.Balance < orderTotalPrice {
		return core.ErrConfict.WithError(orderEntity.ErrInsufficientBalance.Error())
	}

	user.AddBalance(-orderTotalPrice)

	err = uc.repo.CreateOrder(ctx, &newOrder, user, &products)
	if err != nil {
		return core.ErrConfict.WithError(orderEntity.ErrCannotCreateOrder.Error()).WithDebug(err.Error())
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
