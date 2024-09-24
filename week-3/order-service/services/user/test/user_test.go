package test

import (
	"order_service/services/user/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	user entity.User
}

func (suite *UserTestSuite) SetupTest() {
	suite.user = entity.User{
		Id:        1,
		Username:  "partridge",
		Password:  "130703",
		Balance:   0.0,
		CreatedAt: time.Now(),
	}
}

func (suite *UserTestSuite) TestNewUser() {
	user := entity.NewUser(1, "partridge", "130703")

	assert.Equal(suite.T(), 1, user.Id, "Id should be set correctly")
	assert.Equal(suite.T(), "partridge", user.Username, "Username should be set correctly")
	assert.Equal(suite.T(), "130703", user.Password, "Password should be set correctly")
	assert.Equal(suite.T(), float32(0.0), user.Balance, "Balance should be set correctly")
}

func (suite *UserTestSuite) TestSetBalance() {
	tests := []struct {
		name    string
		user    *entity.User
		balance float32
		want    float32
		panic   bool
	}{
		{
			name:    "Non Nil User",
			user:    &suite.user,
			balance: 10.0,
			want:    10.0,
			panic:   false,
		},
		{
			name:    "Nil User",
			user:    nil,
			balance: 10.0,
			want:    0.0,
			panic:   true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.panic {
				assert.NotPanics(suite.T(), func() {
					tt.user.SetBalance(tt.balance)
				}, "Calling SetBalance on nil user should not panic")
			} else {
				tt.user.SetBalance(tt.balance)
				assert.Equal(suite.T(), tt.want, tt.user.Balance, "Balance should be updated")
			}
		})
	}
}

func (suite *UserTestSuite) TestGetId() {
	got := suite.user.GetId()

	assert.Equal(suite.T(), suite.user.Id, got, "Id should be retrieved correctly")
}

func (suite *UserTestSuite) TestGetBalance() {
	got := suite.user.GetBalance()

	assert.Equal(suite.T(), suite.user.Balance, got, "Balance should be retrieved correctly")
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
