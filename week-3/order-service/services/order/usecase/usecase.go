package usecase

import (
	"context"
	"order_service/internal/core"
	"order_service/services/order/entity"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, data *entity.Order) error
}

type orderUsecase struct {
	repo OrderRepo
}

func NewUsecase(repo OrderRepo) *orderUsecase {
	return &orderUsecase{
		repo,
	}
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, data *entity.OrderRequest) error {
	if err := data.Validate(); err != nil {
		return core.ErrBadRequest.WithError(err.Error())
	}

	requester := core.GetRequester(ctx)

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}
	requesterId := int(uid.GetLocalID())

	newOrder := entity.NewOrder(0, requesterId, 0.0, data.Items)

	err = uc.repo.CreateOrder(ctx, &newOrder)
	if err != nil {
		return core.ErrInternalServerError.WithError(err.Error()).WithDebug(err.Error())
	}

	return nil
}
