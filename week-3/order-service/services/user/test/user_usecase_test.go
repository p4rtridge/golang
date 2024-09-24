package test

import (
	"context"
	"errors"
	"order_service/internal/core"
	"order_service/services/user/entity"
	"order_service/services/user/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockUserRepo struct {
	mock.Mock
}

func (mock *mockUserRepo) GetUsers(ctx context.Context) (*[]entity.User, error) {
	args := mock.Called(ctx)

	return args.Get(0).(*[]entity.User), args.Error(1)
}

func (mock *mockUserRepo) GetUserById(ctx context.Context, userId int) (*entity.User, error) {
	args := mock.Called(ctx, userId)

	return args.Get(0).(*entity.User), args.Error(1)
}

func (mock *mockUserRepo) AddUserBalanceById(ctx context.Context, userId int, balance float32) error {
	args := mock.Called(ctx, userId, balance)

	return args.Error(0)
}

type UserUsecaseTestSuite struct {
	suite.Suite
	users        *[]entity.User
	mockUserRepo *mockUserRepo
	usecase      usecase.UserUsecase
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.users = &[]entity.User{
		{
			Id:       1,
			Username: "partridge",
			Password: "130703",
			Balance:  10.0,
		},
		{
			Id:       2,
			Username: "partridge1307",
			Password: "130703",
			Balance:  5.0,
		},
	}

	suite.mockUserRepo = new(mockUserRepo)
	suite.usecase = usecase.NewUsecase(suite.mockUserRepo)
}

func (suite *UserUsecaseTestSuite) TestGetUsers() {
	tests := []struct {
		name      string
		ctx       context.Context
		repoErr   error
		want      *[]entity.User
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Users exists",
			ctx:       context.TODO(),
			repoErr:   nil,
			want:      suite.users,
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "Users empty",
			ctx:       context.TODO(),
			repoErr:   core.ErrRecordNotFound,
			want:      nil,
			wantErr:   core.ErrNotFound,
			assertion: assert.Error,
		},
		{
			name:      "Users return an error",
			ctx:       context.TODO(),
			repoErr:   errors.New("this is an error"),
			want:      nil,
			wantErr:   core.ErrNotFound.WithError(entity.ErrCannotGetUser.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockUserRepo.On("GetUsers", tt.ctx).Return(tt.want, tt.repoErr)

			users, err := suite.usecase.GetUsers(tt.ctx)

			assert.Equal(suite.T(), tt.want, users, "users should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				assert.ErrorIs(suite.T(), err, tt.wantErr, "error should be return correctly")
			}

			suite.mockUserRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserUsecaseTestSuite) TestGetUser() {
	tests := []struct {
		name      string
		ctx       context.Context
		userId    int
		repoErr   error
		want      *entity.User
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "User exists",
			ctx:       context.TODO(),
			userId:    1,
			repoErr:   nil,
			want:      &(*suite.users)[0],
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "User not found",
			ctx:       context.TODO(),
			userId:    1,
			repoErr:   core.ErrRecordNotFound,
			want:      nil,
			wantErr:   core.ErrNotFound,
			assertion: assert.Error,
		},
		{
			name:      "User return an error",
			ctx:       context.TODO(),
			userId:    1,
			repoErr:   errors.New("this is an error"),
			want:      nil,
			wantErr:   core.ErrNotFound.WithError(entity.ErrCannotGetUser.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockUserRepo.On("GetUserById", tt.ctx, tt.userId).Return(tt.want, tt.repoErr)

			user, err := suite.usecase.GetUser(tt.ctx, tt.userId)

			assert.Equal(suite.T(), tt.want, user, "user should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				assert.ErrorIs(suite.T(), err, tt.wantErr, "error should be return correctly")
			}

			suite.mockUserRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserUsecaseTestSuite) TestAddUserBalance() {
	tests := []struct {
		name      string
		ctx       context.Context
		userId    int
		balance   float32
		repoErr   error
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "User exists",
			ctx:       context.TODO(),
			userId:    1,
			balance:   10.0,
			repoErr:   nil,
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name:      "User return an error",
			ctx:       context.TODO(),
			userId:    1,
			balance:   10.0,
			repoErr:   errors.New("this is an error"),
			want:      core.ErrInternalServerError.WithError(entity.ErrCannotAddBalance.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockUserRepo.On("AddUserBalanceById", tt.ctx, tt.userId, tt.balance).Return(tt.repoErr)

			err := suite.usecase.AddUserBalance(tt.ctx, tt.userId, tt.balance)

			if tt.assertion(suite.T(), err) {
				assert.ErrorIs(suite.T(), err, tt.want, "error should be return correctly")
			}

			suite.mockUserRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserUsecaseTestSuite) TestGetUserProfile() {
	tests := []struct {
		name      string
		ctx       context.Context
		userId    int
		repoErr   error
		want      *entity.User
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "User exists",
			ctx:       context.TODO(),
			userId:    1,
			repoErr:   nil,
			want:      &(*suite.users)[0],
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "User return an error",
			ctx:       context.TODO(),
			userId:    1,
			repoErr:   errors.New("this is an error"),
			want:      nil,
			wantErr:   core.ErrUnauthorized.WithError(entity.ErrCannotGetUser.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockUserRepo.On("GetUserById", tt.ctx, tt.userId).Return(tt.want, tt.repoErr)

			user, err := suite.usecase.GetUserProfile(tt.ctx, tt.userId)

			assert.Equal(suite.T(), tt.want, user, "user should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				assert.ErrorIs(suite.T(), err, tt.wantErr, "error should be return correctly")
			}

			suite.mockUserRepo.AssertExpectations(suite.T())
		})
	}
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}
