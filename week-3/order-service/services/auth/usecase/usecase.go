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

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID int, deviceID, token string, expiration int) error
	GetRefreshToken(ctx context.Context, userID int, deviceID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID int, deviceID string) error
	DeleteAllRefreshToken(ctx context.Context, userID int) error
}

type Hasher interface {
	HashPassword(password string) (string, error)
	CompareHash(hashedPassword, password string) (bool, error)
}

type JWT interface {
	IssueAccessToken(ctx context.Context, id, sub string) (string, int, error)
	IssueRefreshToken(ctx context.Context, id, sub string) (string, int, error)
	ParseToken(ctx context.Context, tokenStr string) (*jwt.RegisteredClaims, error)
}

type authUsecase struct {
	repo      AuthRepository
	tokenRepo TokenRepository
	hasher    Hasher
	jwt       JWT
}

func NewUsecase(repo AuthRepository, tokenRepo TokenRepository, hasher Hasher, jwt JWT) *authUsecase {
	return &authUsecase{
		repo,
		tokenRepo,
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

func (uc *authUsecase) Login(ctx context.Context, data *authEntity.AuthLogin) (*authEntity.TokenResponse, error) {
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

func (uc *authUsecase) Refresh(ctx context.Context, data *authEntity.RefreshTokenRequest) (*authEntity.TokenResponse, error) {
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

func (uc *authUsecase) SignOut(ctx context.Context, data *authEntity.AuthSignOut) error {
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
