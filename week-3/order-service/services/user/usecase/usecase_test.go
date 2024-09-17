package usecase_test

import (
	"context"
	"order_service/internal/core"
	"order_service/services/user/entity"
	"order_service/services/user/repository/postgres"
	"order_service/services/user/usecase"
	"testing"

	"github.com/google/uuid"
)

var mockUserData []entity.User = []entity.User{
	{
		Id:       1,
		Username: "partridge",
		Password: "130703",
		Balance:  100.0,
	},
	{
		Id:       2,
		Username: "doggo",
		Password: "130703",
		Balance:  50.0,
	},
	{
		Id:       3,
		Username: "katty",
		Password: "130703",
		Balance:  75.0,
	},
}

type mockRepo struct {
	data []entity.User
}

func NewMockRepo() postgres.UserRepository {
	return &mockRepo{
		data: mockUserData,
	}
}

func (repo *mockRepo) GetUsers(ctx context.Context) (*[]entity.User, error) {
	return &repo.data, nil
}

func (repo *mockRepo) GetUserById(ctx context.Context, userId int) (*entity.User, error) {
	var targetUser *entity.User

	for _, user := range repo.data {
		if user.GetId() == userId {
			targetUser = &user
		}
	}

	if targetUser == nil {
		return nil, entity.ErrCannotGetUser
	}

	return targetUser, nil
}

func (repo *mockRepo) AddUserBalanceById(ctx context.Context, userId int, balance float32) error {
	for idx, user := range repo.data {
		if user.GetId() == userId {
			repo.data[idx].SetBalance(user.GetBalance() + balance)

			return nil
		}
	}

	return entity.ErrCannotAddBalance
}

func TestGetUsers(t *testing.T) {
	uc := usecase.NewUsecase(NewMockRepo())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	users, err := uc.GetUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for idx, user := range *users {
		if user != mockUserData[idx] {
			t.Errorf("expected %v, got %v at index: %d", mockUserData[idx], user, idx)
		}
	}
}

func TestGetUser(t *testing.T) {
	uc := usecase.NewUsecase(NewMockRepo())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetUserId := 1
	user, err := uc.GetUser(ctx, targetUserId)
	if err != nil {
		t.Fatal(err)
	}

	if *user != mockUserData[0] {
		t.Errorf("expected %v, got %v", mockUserData[0], user)
	}
}

func TestGetUserProfile(t *testing.T) {
	uc := usecase.NewUsecase(NewMockRepo())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockUID := core.NewUID(1)
	mockRequester := core.NewRequester(mockUID.String(), uuid.New().String())

	ctx = core.ContextWithRequester(ctx, mockRequester)

	profile, err := uc.GetUserProfile(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if *profile != mockUserData[0] {
		t.Errorf("expected %v, got %v", mockUserData[0], profile)
	}
}

func TestAddUseBalance(t *testing.T) {
	uc := usecase.NewUsecase(NewMockRepo())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockUID := core.NewUID(1)
	mockRequester := core.NewRequester(mockUID.String(), uuid.New().String())

	ctx = core.ContextWithRequester(ctx, mockRequester)

	expectedBalance := mockUserData[0].Balance + 100.0
	err := uc.AddUserBalance(ctx, &entity.UserRequest{
		Balance: 100.0,
	})
	if err != nil {
		t.Fatal(err)
	}

	if mockUserData[0].Balance != expectedBalance {
		t.Errorf("expected %v, got %v", expectedBalance, mockUserData[0].Balance)
	}
}
