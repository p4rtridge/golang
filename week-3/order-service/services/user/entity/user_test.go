package entity_test

import (
	"order_service/services/user/entity"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	mock_user := entity.User{
		Id:        1,
		Username:  "partridge",
		Password:  "130703",
		Balance:   0.0,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	user := entity.NewUser(1, "partridge", "130703")

	if user.Id != mock_user.Id || user.Username != mock_user.Username || user.Password != mock_user.Password || user.Balance != mock_user.Balance {
		t.Errorf("expected %v, got %v", mock_user, user)
	}
}

func TestSetBalance(t *testing.T) {
	expectedBalance := float32(10.0)

	user := entity.NewUser(1, "partridge", "130703")

	user.SetBalance(expectedBalance)

	if user.Balance != expectedBalance {
		t.Errorf("expected %v, got %v", expectedBalance, user.Balance)
	}
}

func TestGetId(t *testing.T) {
	expectedId := 1

	user := entity.NewUser(1, "partridge", "130703")

	if user.GetId() != expectedId {
		t.Errorf("expected %v, got %v", expectedId, user.GetId())
	}
}

func TestGetBalance(t *testing.T) {
	expectedBalance := float32(10.0)

	user := entity.NewUser(1, "partridge", "130703")

	user.SetBalance(expectedBalance)

	if user.GetBalance() != expectedBalance {
		t.Errorf("expected %v, got %v", expectedBalance, user.GetId())
	}
}
