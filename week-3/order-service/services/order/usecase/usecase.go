package usecase

import (
	"context"
	"order_service/internal/core"
	"order_service/services/order/entity"
	orderRepo "order_service/services/order/repository/postgres"
)

type orderUsecase struct {
	repo orderRepo.OrderRepository
}

func NewUsecase(repo orderRepo.OrderRepository) *orderUsecase {
	return &orderUsecase{
		repo,
	}
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, data *entity.OrderRequest) error {
	if err := data.Validate(); err != nil {
		return core.ErrBadRequest.WithError(err.Error())
	}

	// requester := core.GetRequester(ctx)
	//
	// uid, err := core.DecomposeUID(requester.GetSubject())
	// if err != nil {
	// 	return core.ErrInternalServerError.WithDebug(err.Error())
	// }
	// requesterId := int(uid.GetLocalID())

	return nil
}
