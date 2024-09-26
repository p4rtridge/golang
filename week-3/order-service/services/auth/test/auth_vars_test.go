package test

import (
	"order_service/pkg"
	"order_service/services/auth/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthVarsTestSuite struct {
	suite.Suite
}

func (suite *AuthVarsTestSuite) TestAuthUsernamePasswordValidate() {
	tests := []struct {
		name      string
		data      entity.AuthUsernamePassword
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Validated",
			data: entity.AuthUsernamePassword{
				Username: "partridge1307",
				Password: "Duc13072003",
			},
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name: "Invalid username",
			data: entity.AuthUsernamePassword{
				Username: "1307partridge",
				Password: "Duc13072003",
			},
			want:      pkg.ErrUsernameIsNotValid,
			assertion: assert.Error,
		},
		{
			name: "Invalid password",
			data: entity.AuthUsernamePassword{
				Username: "partridge1307",
				Password: "duc13072003",
			},
			want:      pkg.ErrPasswordIsNotValid,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := tt.data.Validate()

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.want, "target value should be valid")
			}
		})
	}
}

func (suite *AuthVarsTestSuite) TestAuthLoginValidate() {
	tests := []struct {
		name      string
		data      entity.AuthLogin
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Validated",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "partridge1307",
					Password: "Duc13072003",
				},
				DeviceId: uuid.New().String(),
			},
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name: "Invalid username",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "1307partridge",
					Password: "Duc13072003",
				},
				DeviceId: uuid.New().String(),
			},
			want:      pkg.ErrUsernameIsNotValid,
			assertion: assert.Error,
		},
		{
			name: "Invalid password",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "partridge1307",
					Password: "duc13072003",
				},
				DeviceId: uuid.New().String(),
			},
			want:      pkg.ErrPasswordIsNotValid,
			assertion: assert.Error,
		},
		{
			name: "Invalid device id",
			data: entity.AuthLogin{
				AuthUsernamePassword: entity.AuthUsernamePassword{
					Username: "partridge1307",
					Password: "Duc13072003",
				},
				DeviceId: "invalid device id",
			},
			want:      pkg.ErrDeviceIdIsNotValid,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := tt.data.Validate()

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.want, "target value should be valid")
			}
		})
	}
}

func (suite *AuthVarsTestSuite) TestAuthSignOutValidate() {
	tests := []struct {
		name      string
		data      entity.AuthSignOut
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Validated",
			data: entity.AuthSignOut{
				DeviceId: uuid.New().String(),
			},
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name: "Invalid device id",
			data: entity.AuthSignOut{
				DeviceId: "invalid device id",
			},
			want:      pkg.ErrDeviceIdIsNotValid,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := tt.data.Validate()

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.want, "target value should be valid")
			}
		})
	}
}

func TestAuthVarsTestSuite(t *testing.T) {
	suite.Run(t, new(AuthVarsTestSuite))
}
