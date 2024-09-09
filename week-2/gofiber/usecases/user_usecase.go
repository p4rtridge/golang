package usecases

import (
	"context"

	"github.com/partridge1307/gofiber/entities"
	"github.com/partridge1307/gofiber/repositories"
)

type UserUsecase interface {
	GetUsers(ctx context.Context) (*[]entities.User, error)
	GetUser(ctx context.Context, u interface{}) (*entities.User, error)
}

type userUsecase struct {
	userRepo repositories.UserRepository
}

func NewUserUsecase(r repositories.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: r,
	}
}

func (uc *userUsecase) GetUsers(ctx context.Context) (*[]entities.User, error) {
	return uc.userRepo.GetUsers(ctx)
}

func (uc *userUsecase) GetUser(ctx context.Context, u interface{}) (*entities.User, error) {
	return uc.userRepo.GetUser(ctx, u)
}
