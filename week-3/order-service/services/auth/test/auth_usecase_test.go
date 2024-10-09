package test

import (
	"context"
	"errors"
	"order_service/internal/core"
	"order_service/services/auth/entity"
	"order_service/services/auth/test/mock"
	"order_service/services/auth/usecase"
	userEntity "order_service/services/user/entity"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type AuthUsecaseTestSuite struct {
	suite.Suite

	mockRepo      *mock.MockAuthRepository
	mockTokenRepo *mock.MockTokenRepository
	mockHasher    *mock.MockHasher
	mockJWT       *mock.MockJWT

	usecase usecase.AuthUseCase
}

func (suite *AuthUsecaseTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.mockRepo = mock.NewMockAuthRepository(ctrl)
	suite.mockTokenRepo = mock.NewMockTokenRepository(ctrl)
	suite.mockHasher = mock.NewMockHasher(ctrl)
	suite.mockJWT = mock.NewMockJWT(ctrl)

	suite.usecase = usecase.NewUsecase(suite.mockRepo, suite.mockTokenRepo, suite.mockHasher, suite.mockJWT)
}

func (suite *AuthUsecaseTestSuite) TestRegister() {
	type authRepo struct {
		err        error
		returnData *userEntity.User
	}

	type hasher struct {
		err        error
		returnData string
	}

	tests := []struct {
		name      string
		data      entity.AuthUsernamePassword
		authRepo  authRepo
		hasher    hasher
		repoErr   error
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Successful register",
			data: entity.AuthUsernamePassword{
				Username: "partridge1307",
				Password: "Duc13072003",
			},
			authRepo: authRepo{
				err:        core.ErrRecordNotFound,
				returnData: nil,
			},
			hasher: hasher{
				err:        nil,
				returnData: "hashed-password",
			},
			repoErr:   nil,
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name: "Username existed",
			data: entity.AuthUsernamePassword{
				Username: "partridge1307",
				Password: "Duc13072003",
			},
			authRepo: authRepo{
				err: nil,
				returnData: &userEntity.User{
					Id:       1,
					Username: "partridge1307",
					Password: "partridge1307",
					Balance:  10,
				},
			},
			hasher: hasher{
				err:        nil,
				returnData: "hashed-password",
			},
			repoErr:   nil,
			want:      core.ErrBadRequest.WithError(entity.ErrUsernameExisted.Error()),
			assertion: assert.Error,
		},
		{
			name: "GetAuth return unexpect error",
			data: entity.AuthUsernamePassword{
				Username: "partridge1307",
				Password: "Duc13072003",
			},
			authRepo: authRepo{
				err: errors.New("this is an error"),
				returnData: &userEntity.User{
					Id:       1,
					Username: "partridge1307",
					Password: "partridge1307",
					Balance:  10,
				},
			},
			hasher: hasher{
				err:        nil,
				returnData: "hashed-password",
			},
			repoErr:   nil,
			want:      core.ErrInternalServerError.WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
		{
			name: "Hash password return an error",
			data: entity.AuthUsernamePassword{
				Username: "partridge1307",
				Password: "Duc13072003",
			},
			authRepo: authRepo{
				err:        core.ErrRecordNotFound,
				returnData: nil,
			},
			hasher: hasher{
				err:        errors.New("unknown errors"),
				returnData: "",
			},
			repoErr:   nil,
			want:      core.ErrInternalServerError.WithError(entity.ErrCannotRegister.Error()).WithDebug(errors.New("unknown errors").Error()),
			assertion: assert.Error,
		},
		{
			name: "AddAuth return an error",
			data: entity.AuthUsernamePassword{
				Username: "partridge1307",
				Password: "Duc13072003",
			},
			authRepo: authRepo{
				err:        core.ErrRecordNotFound,
				returnData: nil,
			},
			hasher: hasher{
				err:        nil,
				returnData: "hashed-password",
			},
			repoErr:   errors.New("this is an error"),
			want:      core.ErrInternalServerError.WithError(entity.ErrCannotRegister.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			authRepo := suite.mockRepo.EXPECT().GetAuth(gomock.Any(), tt.data.Username).Return(tt.authRepo.returnData, tt.authRepo.err)
			hasherRepo := suite.mockHasher.EXPECT().HashPassword(tt.data.Password).AnyTimes().After(authRepo).Return(tt.hasher.returnData, tt.hasher.err)
			suite.mockRepo.EXPECT().AddAuth(gomock.Any(), gomock.Any()).AnyTimes().After(hasherRepo).Return(tt.repoErr)

			err := suite.usecase.Register(context.Background(), tt.data)

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.want, "error should be return correctly")
			}
		})
	}
}

func (suite *AuthUsecaseTestSuite) TestLogin() {
	type authRepo struct {
		err        error
		returnData *userEntity.User
	}

	type tokenRepo struct {
		err error
	}

	type hasher struct {
		err        error
		returnData bool
	}

	type jwt struct {
		token    string
		expireIn int
		err      error
		rtErr    error
	}

	token := entity.Token{
		Token:     "super-security-token",
		ExpiredIn: 60,
	}

	tests := []struct {
		name      string
		data      entity.AuthLogin
		authRepo  authRepo
		tokenRepo tokenRepo
		hasher    hasher
		jwt       jwt
		want      *entity.TokenResponse
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Successful login",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "partridge1307",
					Password: "Duc13072003",
				},
				DeviceId: uuid.New().String(),
			},
			authRepo: authRepo{
				err: nil,
				returnData: &userEntity.User{
					Id:        1,
					Username:  "partridge1307",
					Password:  "hashed-password",
					Balance:   10.0,
					CreatedAt: time.Now(),
				},
			},
			hasher: hasher{
				err:        nil,
				returnData: true,
			},
			jwt: jwt{
				token:    token.Token,
				expireIn: token.ExpiredIn,
				err:      nil,
				rtErr:    nil,
			},
			tokenRepo: tokenRepo{
				err: nil,
			},
			want: &entity.TokenResponse{
				AccessToken:  token,
				RefreshToken: token,
			},
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name: "User not exists",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "partridge1307",
					Password: "Duc13072003",
				},
				DeviceId: uuid.New().String(),
			},
			authRepo: authRepo{
				err: core.ErrRecordNotFound,
				returnData: &userEntity.User{
					Id:        1,
					Username:  "partridge1307",
					Password:  "hashed-password",
					Balance:   10.0,
					CreatedAt: time.Now(),
				},
			},
			hasher: hasher{
				err:        nil,
				returnData: true,
			},
			jwt: jwt{
				token:    token.Token,
				expireIn: token.ExpiredIn,
				err:      nil,
				rtErr:    nil,
			},
			tokenRepo: tokenRepo{
				err: nil,
			},
			want:      nil,
			wantErr:   core.ErrBadRequest.WithError(entity.ErrLoginFailed.Error()),
			assertion: assert.Error,
		},
		{
			name: "GetAuth return an error",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "partridge1307",
					Password: "Duc13072003",
				},
				DeviceId: uuid.New().String(),
			},
			authRepo: authRepo{
				err: errors.New("this is an error"),
				returnData: &userEntity.User{
					Id:        1,
					Username:  "partridge1307",
					Password:  "hashed-password",
					Balance:   10.0,
					CreatedAt: time.Now(),
				},
			},
			hasher: hasher{
				err:        nil,
				returnData: true,
			},
			jwt: jwt{
				token:    token.Token,
				expireIn: token.ExpiredIn,
				err:      nil,
				rtErr:    nil,
			},
			tokenRepo: tokenRepo{
				err: nil,
			},
			want:      nil,
			wantErr:   core.ErrInternalServerError.WithError(entity.ErrLoginFailed.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
		{
			name: "Hasher return an error",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "partridge1307",
					Password: "Duc13072003",
				},
				DeviceId: uuid.New().String(),
			},
			authRepo: authRepo{
				err: nil,
				returnData: &userEntity.User{
					Id:        1,
					Username:  "partridge1307",
					Password:  "hashed-password",
					Balance:   10.0,
					CreatedAt: time.Now(),
				},
			},
			hasher: hasher{
				err:        errors.New("this is an error"),
				returnData: true,
			},
			jwt: jwt{
				token:    token.Token,
				expireIn: token.ExpiredIn,
				err:      nil,
				rtErr:    nil,
			},
			tokenRepo: tokenRepo{
				err: nil,
			},
			want:      nil,
			wantErr:   core.ErrInternalServerError.WithError(entity.ErrLoginFailed.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
		{
			name: "Password does not match",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "partridge1307",
					Password: "Duc13072003",
				},
				DeviceId: uuid.New().String(),
			},
			authRepo: authRepo{
				err: nil,
				returnData: &userEntity.User{
					Id:        1,
					Username:  "partridge1307",
					Password:  "hashed-password",
					Balance:   10.0,
					CreatedAt: time.Now(),
				},
			},
			hasher: hasher{
				err:        nil,
				returnData: false,
			},
			jwt: jwt{
				token:    token.Token,
				expireIn: token.ExpiredIn,
				err:      nil,
				rtErr:    nil,
			},
			tokenRepo: tokenRepo{
				err: nil,
			},
			want:      nil,
			wantErr:   core.ErrBadRequest.WithError(entity.ErrLoginFailed.Error()),
			assertion: assert.Error,
		},
		{
			name: "Access token error",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "partridge1307",
					Password: "Duc13072003",
				},
				DeviceId: uuid.New().String(),
			},
			authRepo: authRepo{
				err: nil,
				returnData: &userEntity.User{
					Id:        1,
					Username:  "partridge1307",
					Password:  "hashed-password",
					Balance:   10.0,
					CreatedAt: time.Now(),
				},
			},
			hasher: hasher{
				err:        nil,
				returnData: true,
			},
			jwt: jwt{
				token:    token.Token,
				expireIn: token.ExpiredIn,
				err:      errors.New("access token error"),
				rtErr:    nil,
			},
			tokenRepo: tokenRepo{
				err: nil,
			},
			want:      nil,
			wantErr:   core.ErrInternalServerError.WithError(entity.ErrLoginFailed.Error()).WithDebug(errors.New("access token error").Error()),
			assertion: assert.Error,
		},
		{
			name: "Refresh token error",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "partridge1307",
					Password: "Duc13072003",
				},
				DeviceId: uuid.New().String(),
			},
			authRepo: authRepo{
				err: nil,
				returnData: &userEntity.User{
					Id:        1,
					Username:  "partridge1307",
					Password:  "hashed-password",
					Balance:   10.0,
					CreatedAt: time.Now(),
				},
			},
			hasher: hasher{
				err:        nil,
				returnData: true,
			},
			jwt: jwt{
				token:    token.Token,
				expireIn: token.ExpiredIn,
				err:      nil,
				rtErr:    errors.New("refresh token error"),
			},
			tokenRepo: tokenRepo{
				err: nil,
			},
			want:      nil,
			wantErr:   core.ErrInternalServerError.WithError(entity.ErrLoginFailed.Error()).WithDebug(errors.New("refresh token error").Error()),
			assertion: assert.Error,
		},
		{
			name: "tokenRepo return an error",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "partridge1307",
					Password: "Duc13072003",
				},
				DeviceId: uuid.New().String(),
			},
			authRepo: authRepo{
				err: nil,
				returnData: &userEntity.User{
					Id:        1,
					Username:  "partridge1307",
					Password:  "hashed-password",
					Balance:   10.0,
					CreatedAt: time.Now(),
				},
			},
			hasher: hasher{
				err:        nil,
				returnData: true,
			},
			jwt: jwt{
				token:    token.Token,
				expireIn: token.ExpiredIn,
				err:      nil,
				rtErr:    nil,
			},
			tokenRepo: tokenRepo{
				err: errors.New("this is an error"),
			},
			want:      nil,
			wantErr:   core.ErrInternalServerError.WithError(entity.ErrLoginFailed.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			authRepo := suite.mockRepo.EXPECT().GetAuth(gomock.Any(), tt.data.Username).Return(tt.authRepo.returnData, tt.authRepo.err)
			hasherRepo := suite.mockHasher.EXPECT().CompareHash(tt.authRepo.returnData.Password, tt.data.Password).AnyTimes().After(authRepo).Return(tt.hasher.returnData, tt.hasher.err)
			atRepo := suite.mockJWT.EXPECT().IssueAccessToken(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().After(hasherRepo).Return(tt.jwt.token, tt.jwt.expireIn, tt.jwt.err)
			rtRepo := suite.mockJWT.EXPECT().IssueRefreshToken(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().After(atRepo).Return(tt.jwt.token, tt.jwt.expireIn, tt.jwt.rtErr)
			suite.mockTokenRepo.EXPECT().SetRefreshToken(gomock.Any(), tt.authRepo.returnData.Id, tt.data.DeviceId, tt.jwt.token, tt.jwt.expireIn).AnyTimes().After(rtRepo).Return(tt.tokenRepo.err)

			token, err := suite.usecase.Login(context.Background(), tt.data)

			suite.Equal(tt.want, token, "token must be set correctly")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error must be return correctly")
			}
		})
	}
}

func (suite *AuthUsecaseTestSuite) TestVerify() {
	tests := []struct {
		name      string
		token     string
		repoErr   error
		want      *jwt.RegisteredClaims
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:    "Validated token",
			token:   "token",
			repoErr: nil,
			want: &jwt.RegisteredClaims{
				Subject: "userId",
				ID:      "unique token id",
			},
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:    "Invalidated token",
			token:   "token",
			repoErr: errors.New("invalid"),
			want: &jwt.RegisteredClaims{
				Subject: "",
				ID:      "",
			},
			wantErr:   core.ErrUnauthorized.WithError(entity.ErrLoginFailed.Error()).WithDebug(errors.New("invalid").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockJWT.EXPECT().ParseToken(gomock.Any(), tt.token).Return(tt.want, tt.repoErr)

			sub, tid, err := suite.usecase.Verify(context.Background(), tt.token)

			suite.Equal(tt.want.Subject, sub, "Sub must be equal")
			suite.Equal(tt.want.ID, tid, "Token ID must be equal")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be returned correctly")
			}
		})
	}
}

func (suite *AuthUsecaseTestSuite) TestSignOutAll() {
	tests := []struct {
		name      string
		ctx       context.Context
		repoErr   error
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Successful sign out",
			ctx:       core.ContextWithRequester(context.Background(), core.NewRequester(core.NewUID(1, 0).String(), "tokenId")),
			repoErr:   nil,
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name:      "Successful sign out",
			ctx:       core.ContextWithRequester(context.Background(), core.NewRequester(core.NewUID(1, 0).String(), "tokenId")),
			repoErr:   nil,
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name:      "Token repo return an error",
			ctx:       core.ContextWithRequester(context.Background(), core.NewRequester(core.NewUID(1, 0).String(), "tokenId")),
			repoErr:   errors.New("this is an error"),
			want:      core.ErrNotFound.WithError(entity.ErrSignoutFailed.Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockTokenRepo.EXPECT().DeleteAllRefreshToken(tt.ctx, gomock.Any()).AnyTimes().Return(tt.repoErr)

			err := suite.usecase.SignOutAll(tt.ctx)

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.want, "error should be returned correctly")
			}
		})
	}
}

func TestAuthUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUsecaseTestSuite))
}
