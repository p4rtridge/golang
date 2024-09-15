package usecase

import (
	"context"
	"order_service/internal/core"
	"order_service/services/user/entity"
)

type UserRepository interface {
	GetUsers(ctx context.Context) (*[]entity.User, error)
	GetUserById(ctx context.Context, userId int) (*entity.User, error)
	AddUserBalanceById(ctx context.Context, userId int, balance float32) error
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

func (uc *userUsecase) GetUsers(ctx context.Context) (*[]entity.User, error) {
	users, err := uc.repo.GetUsers(ctx)
	if err != nil {
		if err == core.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}
		return nil, core.ErrNotFound.WithError(entity.ErrCannotGetUser.Error()).WithDebug(err.Error())
	}

	return users, nil
}

func (uc *userUsecase) GetUser(ctx context.Context, userID int) (*entity.User, error) {
	user, err := uc.repo.GetUserById(ctx, userID)
	if err != nil {
		if err == core.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}
		return nil, core.ErrNotFound.WithError(entity.ErrCannotGetUser.Error()).WithDebug(err.Error())
	}

	return user, nil
}

func (uc *userUsecase) AddUserBalance(ctx context.Context, data *entity.UserRequest) error {
	requester := core.GetRequester(ctx)

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}
	requesterId := int(uid.GetLocalID())

	err = uc.repo.AddUserBalanceById(ctx, requesterId, data.Balance)
	if err != nil {
		return core.ErrInternalServerError.WithError(entity.ErrCannotAddBalance.Error()).WithDebug(err.Error())
	}

	return nil
}
