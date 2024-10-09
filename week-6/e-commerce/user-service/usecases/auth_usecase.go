package usecases

import (
	"context"
	"user-service/entities"
	"user-service/helpers"
	"user-service/repositories"

	"github.com/google/uuid"
	"github.com/partridge1307/service-ctx/core"
)

type AuthUsecase interface {
	Register(ctx context.Context, data entities.AuthUsernamePassword) error
	Login(ctx context.Context, data entities.AuthLogin) (*entities.TokenResponse, error)
	Verify(ctx context.Context, token string) (*helpers.ParsedToken, error)
	Refresh(ctx context.Context, data entities.RefreshTokenRequest) (*entities.TokenResponse, error)
	SignOut(ctx context.Context, data entities.AuthSignOut) error
	SignOutAll(ctx context.Context) error
}

type authUsecase struct {
	authRepo  repositories.AuthRepository
	tokenRepo repositories.TokenRepository
	hasher    helpers.Hasher
	jwt       helpers.JWT
}

func NewAuthUsecase(authRepo repositories.AuthRepository, tokenRepo repositories.TokenRepository, hasher helpers.Hasher, jwt helpers.JWT) AuthUsecase {
	return &authUsecase{
		authRepo,
		tokenRepo,
		hasher,
		jwt,
	}
}

func (uc *authUsecase) Register(ctx context.Context, data entities.AuthUsernamePassword) error {
	authData, err := uc.authRepo.GetAuth(ctx, data.Username)

	if err == nil && authData != nil {
		return core.ErrBadRequest.WithError(entities.ErrUsernameExisted.Error())
	}
	if err != nil && err != core.ErrRecordNotFound {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}

	hashedPassword, err := uc.hasher.HashPassword(data.Password)
	if err != nil {
		return core.ErrInternalServerError.WithError(entities.ErrCannotRegister.Error()).WithDebug(err.Error())
	}

	newAuth := entities.NewAuth(data.Username, hashedPassword, 0)

	if err := uc.authRepo.AddAuth(ctx, newAuth); err != nil {
		return core.ErrInternalServerError.WithError(entities.ErrCannotRegister.Error()).WithDebug(err.Error())
	}

	return nil
}

func (uc *authUsecase) Login(ctx context.Context, data entities.AuthLogin) (*entities.TokenResponse, error) {
	authData, err := uc.authRepo.GetAuth(ctx, data.Username)
	if err != nil {
		if err == core.ErrRecordNotFound {
			return nil, core.ErrBadRequest.WithError(entities.ErrLoginFailed.Error())
		}
		return nil, core.ErrInternalServerError.WithError(entities.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	matched, err := uc.hasher.CompareHash(authData.Password, data.Password)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(entities.ErrLoginFailed.Error()).WithDebug(err.Error())
	}
	if !matched {
		return nil, core.ErrBadRequest.WithError(entities.ErrLoginFailed.Error())
	}

	sub := core.NewUID(uint32(authData.Id), uint32(authData.Role)).String()
	tid := uuid.New().String()

	accessTokenStr, atExpireIn, err := uc.jwt.IssueAccessToken(ctx, sub, tid)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(entities.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	refreshTokenStr, rtExpireIn, err := uc.jwt.IssueRefreshToken(ctx, sub, tid)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(entities.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	err = uc.tokenRepo.SetRefreshToken(ctx, authData.Id, data.DeviceId, refreshTokenStr, rtExpireIn)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(entities.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	return &entities.TokenResponse{
		AccessToken: entities.Token{
			Token:     accessTokenStr,
			ExpiredIn: atExpireIn,
		},
		RefreshToken: entities.Token{
			Token:     refreshTokenStr,
			ExpiredIn: rtExpireIn,
		},
	}, nil
}

func (uc *authUsecase) Verify(ctx context.Context, token string) (*helpers.ParsedToken, error) {
	claims, err := uc.jwt.ParseToken(ctx, token)
	if err != nil {
		return nil, core.ErrUnauthorized.WithError(entities.ErrLoginFailed.Error()).WithDebug(err.Error())
	}

	return claims, nil
}

func (uc *authUsecase) Refresh(ctx context.Context, data entities.RefreshTokenRequest) (*entities.TokenResponse, error) {
	claims, err := uc.jwt.ParseToken(ctx, data.RefreshToken)
	if err != nil {
		return nil, core.ErrBadRequest.WithError(entities.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}

	uid, err := core.DecomposeUID(claims.GetSubject())
	if err != nil {
		return nil, core.ErrBadRequest.WithError(entities.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}
	userId := int(uid.GetLocalID())

	dataToken, err := uc.tokenRepo.GetRefreshToken(ctx, userId, data.DeviceId)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(entities.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}

	if dataToken != data.RefreshToken {
		return nil, core.ErrUnauthorized.WithError(entities.ErrRefreshFailed.Error())
	}

	sub := uid.String()
	tid := uuid.New().String()

	accessTokenStr, atExpireIn, err := uc.jwt.IssueAccessToken(ctx, tid, sub)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(entities.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}

	refreshTokenStr, rtExpireIn, err := uc.jwt.IssueRefreshToken(ctx, tid, sub)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(entities.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}

	err = uc.tokenRepo.SetRefreshToken(ctx, userId, data.DeviceId, refreshTokenStr, rtExpireIn)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(entities.ErrRefreshFailed.Error()).WithDebug(err.Error())
	}

	return &entities.TokenResponse{
		AccessToken: entities.Token{
			Token:     accessTokenStr,
			ExpiredIn: atExpireIn,
		},
		RefreshToken: entities.Token{
			Token:     refreshTokenStr,
			ExpiredIn: rtExpireIn,
		},
	}, nil
}

func (uc *authUsecase) SignOut(ctx context.Context, data entities.AuthSignOut) error {
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
		return core.ErrNotFound.WithError(entities.ErrSignoutFailed.Error())
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
		return core.ErrNotFound.WithError(entities.ErrSignoutFailed.Error())
	}

	return nil
}
