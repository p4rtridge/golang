package usecase

import (
	"context"
	"order_service/internal/core"
	"order_service/services/user/entity"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id int) (*entity.User, error)
}

type userUsecase struct {
	repo UserRepository
}

func NewUsecase(repo UserRepository) *userUsecase {
	return &userUsecase{
		repo,
	}
}

func (uc *userUsecase) GetUserProfile(ctx context.Context) (*entity.User, error) {
	requester := core.GetRequester(ctx)

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return nil, core.ErrInternalServerError.WithDebug(err.Error())
	}
	requesterId := int(uid.GetLocalID())

	user, err := uc.repo.GetUserById(ctx, requesterId)
	if err != nil {
		return nil, core.ErrUnauthorized.WithError(entity.ErrCannotGetUser.Error()).WithDebug(err.Error())
	}

	return user, nil
}
