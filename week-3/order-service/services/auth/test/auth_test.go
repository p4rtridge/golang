package test

import (
	"order_service/services/auth/entity"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	auth entity.Auth
}

func (suite *AuthTestSuite) SetupTest() {
	suite.auth = entity.Auth{
		Username: "partridge",
		Password: "130703",
	}
}

func (suite *AuthTestSuite) TestNewAuth() {
	auth := entity.NewAuth("partridge", "130703", 0)

	suite.Equal("partridge", auth.Username, "Username should be set correctly")
	suite.Equal("130703", auth.Password, "Password should be set correctly")
	suite.Equal(0, auth.Role, "Password should be set correctly")
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
