package repositories

import (
	"context"

	"github.com/partridge1307/gofiber/entities"
)

type UserRepository interface {
	GetUsers(ctx context.Context) (*[]entities.User, error)
	GetUser(ctx context.Context, u interface{}) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) error
}
