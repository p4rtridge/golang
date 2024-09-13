package usecase

import (
	"context"
	"order_service/internal/core"
	authEntity "order_service/services/auth/entity"
	userEntity "order_service/services/user/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthRepository interface {
	AddAuth(ctx context.Context, data *authEntity.Auth) error
	GetAuth(ctx context.Context, username string) (*userEntity.User, error)
}

type Hasher interface {
	HashPassword(password string) (string, error)
	CompareHash(hashedPassword, password string) (bool, error)
}

type JWT interface {
	IssueToken(ctx context.Context, id, sub string) (string, int, error)
	ParseToken(ctx context.Context, tokenStr string) (*jwt.RegisteredClaims, error)
}

type authUsecase struct {
	repo   AuthRepository
	hasher Hasher
	jwt    JWT
}

func NewUsecase(repo AuthRepository, hasher Hasher, jwt JWT) *authUsecase {
	return &authUsecase{
		repo,
		hasher,
		jwt,
	}
}

func (uc *authUsecase) Register(ctx context.Context, data *authEntity.AuthUsernamePassword) error {
	if err := data.Validate(); err != nil {
		return core.ErrBadRequest.WithError(err.Error())
	}

	authData, err := uc.repo.GetAuth(ctx, data.Username)

	if err == nil && authData != nil {
		return core.ErrBadRequest.WithError(authEntity.ErrUsernameExisted.Error())
	}
	if err != nil && err != core.ErrRecordNotFound {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}

	hashedPassword, err := uc.hasher.HashPassword(data.Password)
	if err != nil {
		return core.ErrInternalServerError.WithError(authEntity.ErrCannotRegister.Error()).WithDebug(err.Error())
	}

	newAuth := authEntity.NewAuth(data.Username, hashedPassword)

	if err := uc.repo.AddAuth(ctx, &newAuth); err != nil {
		return core.ErrInternalServerError.WithError(authEntity.ErrCannotRegister.Error()).WithDebug(err.Error())
	}

	return nil
}

func (uc *authUsecase) Login(ctx context.Context, data *authEntity.AuthUsernamePassword) (*authEntity.TokenResponse, error) {
	if err := data.Validate(); err != nil {
		return nil, core.ErrBadRequest.WithError(err.Error())
	}

	authData, err := uc.repo.GetAuth(ctx, data.Username)
	if err != nil {
		if err == core.ErrRecordNotFound {
			return nil, core.ErrBadRequest.WithError(authEntity.ErrLoginFailed.Error())
		}
		return nil, core.ErrInternalServerError.WithError(authEntity.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	matched, err := uc.hasher.CompareHash(authData.Password, data.Password)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(authEntity.ErrLoginFailed.Error()).WithDebug(err.Error())
	}
	if !matched {
		return nil, core.ErrBadRequest.WithError(authEntity.ErrLoginFailed.Error())
	}

	uid := core.NewUID(uint32(authData.Id))
	sub := uid.String()
	tid := uuid.New().String()

	tokenStr, expiredIn, err := uc.jwt.IssueToken(ctx, tid, sub)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(authEntity.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	return &authEntity.TokenResponse{
		AccessToken: authEntity.Token{
			Token:     tokenStr,
			ExpiredIn: expiredIn,
		},
	}, nil
}

func (uc *authUsecase) Verify(ctx context.Context, accessToken string) (string, string, error) {
	claims, err := uc.jwt.ParseToken(ctx, accessToken)
	if err != nil {
		return "", "", core.ErrUnauthorized.WithError(authEntity.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	return claims.Subject, claims.ID, nil
}
