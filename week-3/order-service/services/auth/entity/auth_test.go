package entity_test

import (
	"order_service/services/auth/entity"
	"testing"
)

func TestNewAuth(t *testing.T) {
	mock_auth := entity.Auth{
		Username: "partridge",
		Password: "130703",
	}

	auth := entity.NewAuth("partridge", "130703")

	if auth != mock_auth {
		t.Errorf("expected %v, got %v", mock_auth, auth)
	}
}
