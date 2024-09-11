package user

import (
	"context"

	"github.com/partridge1307/gofiber/entity"
)

type UserUsecase interface {
	GetUsers(ctx context.Context) (*[]entity.User, error)
	GetUser(ctx context.Context, u interface{}) (*entity.User, error)
}

type Service struct {
	repo entity.UserRepo
}

func NewService(repo entity.UserRepo) UserUsecase {
	return &Service{
		repo,
	}
}

func (s *Service) GetUsers(ctx context.Context) (*[]entity.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *Service) GetUser(ctx context.Context, u interface{}) (*entity.User, error) {
	user, err := s.repo.GetUser(ctx, u)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &UserNotFoundError{}
	}

	return user, nil
}
