package usecase

import (
	"context"
	"order_service/internal/core"
	"order_service/pkg"
	authEntity "order_service/services/auth/entity"
	authRepo "order_service/services/auth/repository/postgres"
	tokenRepo "order_service/services/auth/repository/redis"

	"github.com/google/uuid"
)

type AuthUseCase interface {
	Register(ctx context.Context, data authEntity.AuthUsernamePassword) error
	Login(ctx context.Context, data authEntity.AuthLogin) (*authEntity.TokenResponse, error)
	Verify(ctx context.Context, token string) (string, string, error)
	Refresh(ctx context.Context, data authEntity.RefreshTokenRequest) (*authEntity.TokenResponse, error)
	SignOut(ctx context.Context, data authEntity.AuthSignOut) error
	SignOutAll(ctx context.Context) error
}

type authUsecase struct {
	repo      authRepo.AuthRepository
	tokenRepo tokenRepo.TokenRepository
	hasher    pkg.Hasher
	jwt       pkg.JWT
}

func NewUsecase(repo authRepo.AuthRepository, tokenRepo tokenRepo.TokenRepository, hasher pkg.Hasher, jwt pkg.JWT) AuthUseCase {
	return &authUsecase{
		repo,
		tokenRepo,
		hasher,
		jwt,
	}
}

func (uc *authUsecase) Register(ctx context.Context, data authEntity.AuthUsernamePassword) error {
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

	newAuth := authEntity.NewAuth(data.Username, hashedPassword, 0)

	if err := uc.repo.AddAuth(ctx, newAuth); err != nil {
		return core.ErrInternalServerError.WithError(authEntity.ErrCannotRegister.Error()).WithDebug(err.Error())
	}

	return nil
}

func (uc *authUsecase) Login(ctx context.Context, data authEntity.AuthLogin) (*authEntity.TokenResponse, error) {
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

	uid := core.NewUID(uint32(authData.Id), uint32(authData.Role))
	sub := uid.String()
	tid := uuid.New().String()

	accessTokenStr, atExpireIn, err := uc.jwt.IssueAccessToken(ctx, tid, sub)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(authEntity.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	refreshTokenStr, rtExpireIn, err := uc.jwt.IssueRefreshToken(ctx, tid, sub)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(authEntity.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	err = uc.tokenRepo.SetRefreshToken(ctx, authData.Id, data.DeviceId, refreshTokenStr, rtExpireIn)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(authEntity.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	return &authEntity.TokenResponse{
		AccessToken: authEntity.Token{
			Token:     accessTokenStr,
			ExpiredIn: atExpireIn,
		},
		RefreshToken: authEntity.Token{
			Token:     refreshTokenStr,
			ExpiredIn: rtExpireIn,
		},
	}, nil
}

func (uc *authUsecase) Verify(ctx context.Context, token string) (string, string, error) {
	claims, err := uc.jwt.ParseToken(ctx, token)
	if err != nil {
		return "", "", core.ErrUnauthorized.WithError(authEntity.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	return claims.Subject, claims.ID, nil
}

func (uc *authUsecase) Refresh(ctx context.Context, data authEntity.RefreshTokenRequest) (*authEntity.TokenResponse, error) {
	claims, err := uc.jwt.ParseToken(ctx, data.RefreshToken)
	if err != nil {
		return nil, core.ErrBadRequest.WithError(authEntity.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}

	uid, err := core.DecomposeUID(claims.Subject)
	if err != nil {
		return nil, core.ErrBadRequest.WithError(authEntity.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}

	userId := int(uid.GetLocalID())

	dataToken, err := uc.tokenRepo.GetRefreshToken(ctx, userId, data.DeviceId)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(authEntity.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}

	if dataToken != data.RefreshToken {
		return nil, core.ErrUnauthorized.WithError(authEntity.ErrRefreshFailed.Error())
	}

	sub := uid.String()
	tid := uuid.New().String()

	accessTokenStr, atExpireIn, err := uc.jwt.IssueAccessToken(ctx, tid, sub)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(authEntity.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}

	refreshTokenStr, rtExpireIn, err := uc.jwt.IssueRefreshToken(ctx, tid, sub)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(authEntity.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}

	err = uc.tokenRepo.SetRefreshToken(ctx, userId, data.DeviceId, refreshTokenStr, rtExpireIn)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(authEntity.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}

	return &authEntity.TokenResponse{
		AccessToken: authEntity.Token{
			Token:     accessTokenStr,
			ExpiredIn: atExpireIn,
		},
		RefreshToken: authEntity.Token{
			Token:     refreshTokenStr,
			ExpiredIn: rtExpireIn,
		},
	}, nil
}

func (uc *authUsecase) SignOut(ctx context.Context, data authEntity.AuthSignOut) error {
	if err := data.Validate(); err != nil {
		return core.ErrUnauthorized.WithError(err.Error())
	}

	requester := core.GetRequester(ctx)

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}

	requesterId := int(uid.GetLocalID())

	err = uc.tokenRepo.DeleteRefreshToken(ctx, requesterId, data.DeviceId)
	if err != nil {
		return core.ErrNotFound.WithError(authEntity.ErrSignoutFailed.Error())
	}

	return nil
}

func (uc *authUsecase) SignOutAll(ctx context.Context) error {
	requester := core.GetRequester(ctx)

	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}

	requesterId := int(uid.GetLocalID())

	err = uc.tokenRepo.DeleteAllRefreshToken(ctx, requesterId)
	if err != nil {
		return core.ErrNotFound.WithError(authEntity.ErrSignoutFailed.Error())
	}

	return nil
}
