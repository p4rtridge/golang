package usecase_test

import (
	"context"
	"errors"
	"order_service/internal/core"
	"order_service/services/order/entity"
	"order_service/services/order/repository/postgres"
	"order_service/services/order/usecase"
	productEntity "order_service/services/product/entity"
	userEntity "order_service/services/user/entity"
	"testing"
	"time"

	"github.com/google/uuid"
)

var mockUserData []userEntity.User = []userEntity.User{
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

var mockProductData []productEntity.Product = []productEntity.Product{
	{
		Id:        1,
		Name:      "your mom",
		Quantity:  4,
		Price:     25.0,
		CreatedAt: time.Now(),
	},
	{
		Id:        2,
		Name:      "your dad",
		Quantity:  5,
		Price:     25.0,
		CreatedAt: time.Now(),
	},
	{
		Id:        3,
		Name:      "your sister",
		Quantity:  1,
		Price:     300.0,
		CreatedAt: time.Now(),
	},
}

var mockOrderData []entity.Order = []entity.Order{
	{
		Id:         1,
		UserId:     1,
		TotalPrice: 100.0,
		Items: []entity.OrderItem{
			entity.NewOrderItem(1, 1, "your mom", 25.0, 2),
			entity.NewOrderItem(1, 2, "your dad", 25.0, 2),
		},
	},
}

type mockRepo struct {
	data        []entity.Order
	userData    []userEntity.User
	productData []productEntity.Product
}

func NewMockRepo() postgres.OrderRepository {
	return &mockRepo{
		data:        mockOrderData,
		userData:    mockUserData,
		productData: mockProductData,
	}
}

func (repo *mockRepo) GetOrders(ctx context.Context, userId int) (*[]entity.Order, error) {
	orders := make([]entity.Order, 0)

	for _, order := range repo.data {
		if order.UserId == userId {
			orders = append(orders, order)
		}
	}

	return &orders, nil
}

func (repo *mockRepo) CreateOrder(ctx context.Context, order *entity.Order, callbackFn func(order *entity.Order, user *userEntity.User, products *[]productEntity.Product) (bool, error)) error {
	var userData *userEntity.User

	for _, user := range repo.userData {
		if order.GetUserIdSafe() == user.GetId() {
			userData = &user
		}
	}

	if userData == nil {
		return errors.New("not found user")
	}

	products := make([]productEntity.Product, 0, len(order.Items))
	for _, order := range order.Items {
		for _, product := range repo.productData {
			if order.GetProductId() == product.Id {
				products = append(products, productEntity.NewProduct(product.Id, product.GetName(), product.GetQuantity(), product.GetPrice()))
			}
		}
	}

	accept, err := callbackFn(order, userData, &products)
	if err != nil {
		return err
	}
	if !accept {
		return entity.ErrCannotCreateOrder
	}

	newOrder := entity.NewOrder(repo.data[len(repo.data)-1].Id+1, order.UserId, order.TotalPrice, order.Items)
	repo.data = append(repo.data, newOrder)

	for _, product := range products {
		for idx, productData := range repo.productData {
			if product.Id == productData.Id {
				repo.productData[idx].SetQuantity(product.GetQuantity())
			}
		}
	}

	for idx, user := range repo.userData {
		if user.Id == userData.Id {
			repo.userData[idx].SetBalance(userData.GetBalance())
		}
	}

	return nil
}

func TestGetOrders(t *testing.T) {
	uc := usecase.NewUsecase(NewMockRepo())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockUID := core.NewUID(1)
	mockRequester := core.NewRequester(mockUID.String(), uuid.New().String())

	ctx = core.ContextWithRequester(ctx, mockRequester)

	orders, err := uc.GetOrders(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(*orders) != len(mockOrderData) {
		t.Errorf("expected %d orders, got %d orders", len(mockOrderData), len(*orders))
	}

	for idx, order := range *orders {
		if order.Id != mockOrderData[idx].Id || order.UserId != mockOrderData[idx].UserId || order.TotalPrice != mockOrderData[idx].TotalPrice {
			t.Errorf("expected %v, got %v", mockOrderData[idx], order)
		}

		for i, item := range order.Items {
			if item != mockOrderData[idx].Items[i] {
				t.Errorf("expected order's item %v, got %v", mockOrderData[idx].Items[i], item)
			}
		}
	}
}

func TestCreateOrder(t *testing.T) {
	uc := usecase.NewUsecase(NewMockRepo())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expectedUserId := 1
	mockUID := core.NewUID(uint32(expectedUserId))
	mockRequester := core.NewRequester(mockUID.String(), uuid.New().String())

	ctx = core.ContextWithRequester(ctx, mockRequester)

	expectedProduct := mockProductData[0]
	expectedQuantity := 1

	newOrder := entity.NewOrder(0, 0, 0.0, []entity.OrderItem{
		{
			ProductId: expectedProduct.Id,
			Quantity:  expectedQuantity,
		},
	})

	err := uc.CreateOrder(ctx, &newOrder)
	if err != nil {
		t.Fatal(err)
	}

	orders, err := uc.GetOrders(ctx)
	if err != nil {
		t.Fatal(err)
	}

	latestOrder := (*orders)[len(*orders)-1]
	if latestOrder.UserId != expectedUserId || latestOrder.TotalPrice != (expectedProduct.GetPrice()*float32(expectedQuantity)) {
		t.Errorf("expected %v, got %v", newOrder, latestOrder)
	}
}
