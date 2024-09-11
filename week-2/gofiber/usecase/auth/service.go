package auth

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/partridge1307/gofiber/entity"
	"github.com/partridge1307/gofiber/pkg/jwt"
	"github.com/partridge1307/gofiber/pkg/password"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, user *entity.User) error
	SignIn(ctx context.Context, user *entity.User) (string, error)
	Verify(ctx context.Context, tokenStr string) (bool, error)
}

type Service struct {
	validator *validator.Validate
	repo      entity.UserRepo
}

func NewService(validator *validator.Validate, repo entity.UserRepo) AuthUsecase {
	return &Service{
		validator,
		repo,
	}
}

func (s *Service) SignUp(ctx context.Context, user *entity.User) error {
	validationOutcome := user.Validate(s.validator)
	if !validationOutcome.Valid {
		return &ValidationError{
			Errors: validationOutcome.Errors,
		}
	}

	existUser, err := s.repo.GetUser(ctx, user.Username)
	if err != nil {
		return err
	}
	if existUser != nil {
		return &UserExistError{}
	}

	hashed, err := password.Hash(user.Password)
	if err != nil {
		return err
	}

	user.SetPassword(hashed)

	return s.repo.CreateUser(ctx, user)
}

func (s *Service) SignIn(ctx context.Context, user *entity.User) (string, error) {
	validationOutcome := user.Validate(s.validator)
	if !validationOutcome.Valid {
		return "", &ValidationError{
			Errors: validationOutcome.Errors,
		}
	}

	dbUser, err := s.repo.GetUser(ctx, user.Username)
	if err != nil {
		return "", err
	}

	matched, err := password.Verify(user.Password, dbUser.Password)
	if err != nil {
		return "", err
	}

	if !matched {
		return "", &UserUnmatchError{}
	}

	return jwt.Sign(dbUser)
}

func (s *Service) Verify(ctx context.Context, tokenStr string) (bool, error) {
	user, err := jwt.Verify(tokenStr)
	if err != nil {
		return false, err
	}

	_, err = s.repo.GetUser(ctx, user.Username)
	if err != nil {
		return false, err
	}

	return true, nil
}
