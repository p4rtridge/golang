package usecases

import (
	"context"
	"errors"

	"github.com/partridge1307/gofiber/entities"
	"github.com/partridge1307/gofiber/helpers"
	"github.com/partridge1307/gofiber/repositories"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, user *entities.User) error
	SignIn(ctx context.Context, user *entities.User) (string, error)
}

type authUsecase struct {
	repo repositories.UserRepository
}

func NewAuthUsecase(repo repositories.UserRepository) AuthUsecase {
	return &authUsecase{
		repo,
	}
}

func (uc *authUsecase) SignUp(ctx context.Context, user *entities.User) error {
	// Whether the usename already exists
	existUser, err := uc.repo.GetUser(ctx, user.Username)
	if err == nil {
		return errors.New("username already exists")
	}
	if existUser == nil {
		return err
	}

	// Hash the password
	hashed, err := helpers.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Set user's password
	user.SetPassword(hashed)

	// Save and return
	return uc.repo.CreateUser(ctx, user)
}

func (uc *authUsecase) SignIn(ctx context.Context, user *entities.User) (string, error) {
	dbUser, err := uc.repo.GetUser(ctx, user.Username)
	if err != nil {
		return "", errors.New("unmatched")
	}

	matched, err := helpers.VerifyPassword(user.Password, dbUser.Password)
	if err != nil {
		return "", err
	}

	if !matched {
		return "nil", errors.New("unmatched")
	}

	token, err := helpers.SignJWT(dbUser)
	if err != nil {
		return "", err
	}

	return token, nil
}
