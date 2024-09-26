package test

import (
	"context"
	"errors"
	"order_service/internal/core"
	"order_service/services/user/entity"
	"order_service/services/user/test/mock"
	"order_service/services/user/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	users    *[]entity.User
	mockRepo *mock.MockUserRepository
	usecase  usecase.UserUsecase
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

	ctrl := gomock.NewController(suite.T())

	suite.mockRepo = mock.NewMockUserRepository(ctrl)
	suite.usecase = usecase.NewUsecase(suite.mockRepo)
}

func (suite *UserUsecaseTestSuite) TestGetUsers() {
	tests := []struct {
		name      string
		repoErr   error
		want      *[]entity.User
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Users exists",
			repoErr:   nil,
			want:      suite.users,
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "Users empty",
			repoErr:   core.ErrRecordNotFound,
			want:      nil,
			wantErr:   core.ErrNotFound,
			assertion: assert.Error,
		},
		{
			name:      "Users return an error",
			repoErr:   errors.New("this is an error"),
			want:      nil,
			wantErr:   core.ErrNotFound.WithError(entity.ErrCannotGetUser.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().GetUsers(gomock.Any()).Return(tt.want, tt.repoErr)

			users, err := suite.usecase.GetUsers(context.Background())

			suite.Equal(tt.want, users, "users should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be return correctly")
			}
		})
	}
}

func (suite *UserUsecaseTestSuite) TestGetUser() {
	tests := []struct {
		name      string
		userId    int
		repoErr   error
		want      *entity.User
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "User exists",
			userId:    1,
			repoErr:   nil,
			want:      &(*suite.users)[0],
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "User not found",
			userId:    1,
			repoErr:   core.ErrRecordNotFound,
			want:      nil,
			wantErr:   core.ErrNotFound,
			assertion: assert.Error,
		},
		{
			name:      "User return an error",
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

			suite.mockRepo.EXPECT().GetUserById(gomock.Any(), tt.userId).Return(tt.want, tt.repoErr)

			user, err := suite.usecase.GetUser(context.Background(), tt.userId)

			suite.Equal(tt.want, user, "user should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be return correctly")
			}
		})
	}
}

func (suite *UserUsecaseTestSuite) TestAddUserBalance() {
	tests := []struct {
		name      string
		userId    int
		balance   float32
		repoErr   error
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "User exists",
			userId:    1,
			balance:   10.0,
			repoErr:   nil,
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name:      "User return an error",
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

			suite.mockRepo.EXPECT().AddUserBalanceById(gomock.Any(), tt.userId, tt.balance).Return(tt.repoErr)

			err := suite.usecase.AddUserBalance(context.Background(), tt.userId, tt.balance)

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.want, "error should be return correctly")
			}
		})
	}
}

func (suite *UserUsecaseTestSuite) TestGetUserProfile() {
	tests := []struct {
		name      string
		userId    int
		repoErr   error
		want      *entity.User
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "User exists",
			userId:    1,
			repoErr:   nil,
			want:      &(*suite.users)[0],
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "User return an error",
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

			suite.mockRepo.EXPECT().GetUserById(gomock.Any(), tt.userId).Return(tt.want, tt.repoErr)

			user, err := suite.usecase.GetUserProfile(context.Background(), tt.userId)

			suite.Equal(tt.want, user, "user should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be return correctly")
			}
		})
	}
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}
