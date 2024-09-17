package usecase_test

import (
	"context"
	"order_service/services/product/entity"
	"order_service/services/product/repository/postgres"
	"order_service/services/product/usecase"
	"testing"
	"time"
)

var mockProductData []entity.Product = []entity.Product{
	{
		Id:        1,
		Name:      "your mom",
		Quantity:  1,
		Price:     100.0,
		CreatedAt: time.Now(),
	},
	{
		Id:        2,
		Name:      "your dad",
		Quantity:  1,
		Price:     200.0,
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

type mockRepo struct {
	Data []entity.Product
}

func NewMockRepo() postgres.ProductRepository {
	return &mockRepo{
		Data: mockProductData,
	}
}

func (repo *mockRepo) CreateProduct(ctx context.Context, data *entity.Product) error {
	repo.Data = append(repo.Data, *data)

	return nil
}

func (repo *mockRepo) GetProducts(ctx context.Context) (*[]entity.Product, error) {
	return &repo.Data, nil
}

func (repo *mockRepo) GetProduct(ctx context.Context, productID int) (*entity.Product, error) {
	var targetProduct *entity.Product

	for _, product := range repo.Data {
		if product.Id == productID {
			targetProduct = &product
		}
	}

	if targetProduct == nil {
		return nil, entity.ErrProductNotFound
	}

	return targetProduct, nil
}

func (repo *mockRepo) UpdateProduct(ctx context.Context, productID int, data *entity.Product) error {
	for idx := range repo.Data {
		currentProduct := &repo.Data[idx]

		if currentProduct.Id == productID {
			now := time.Now()

			currentProduct.Name = data.Name
			currentProduct.Price = data.Price
			currentProduct.SetQuantity(data.Quantity)
			currentProduct.UpdatedAt = &now

			return nil
		}
	}

	return entity.ErrCannotUpdate
}

func (repo *mockRepo) DeleteProduct(ctx context.Context, productID int) error {
	for idx, product := range repo.Data {
		if product.Id == productID {
			repo.Data = append(repo.Data[:idx], repo.Data[idx+1:]...)

			return nil
		}
	}

	return entity.ErrCannotDelete
}

func TestGetProducts(t *testing.T) {
	uc := usecase.NewUsecase(NewMockRepo())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	products, err := uc.GetProducts(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for idx, product := range *products {
		if product != mockProductData[idx] {
			t.Errorf("expected %v, got %v at index %d", mockProductData[idx], product, idx)
		}
	}
}

func TestGetProduct(t *testing.T) {
	uc := usecase.NewUsecase(NewMockRepo())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetProductId := 1
	product, err := uc.GetProduct(ctx, targetProductId)
	if err != nil {
		t.Fatal(err)
	}

	if *product != mockProductData[0] {
		t.Errorf("expected %v, got %v", mockProductData[0], product)
	}
}

func TestCreateProduct(t *testing.T) {
	repo := NewMockRepo()
	uc := usecase.NewUsecase(repo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataReq := entity.ProductRequest{
		Name:     "your step sister",
		Quantity: 1,
		Price:    250.0,
	}

	err := uc.CreateProduct(ctx, &dataReq)
	if err != nil {
		t.Fatal(err)
	}

	products, err := uc.GetProducts(ctx)
	if err != nil {
		t.Fatal(err)
	}

	productLen := len(*products) - 1
	if (*products)[productLen].Name != dataReq.Name || (*products)[productLen].Quantity != dataReq.Quantity || (*products)[productLen].Price != dataReq.Price {
		t.Errorf("expected %v, got %v", dataReq, (*products)[productLen])
	}
}

func TestUpdateProduct(t *testing.T) {
	uc := usecase.NewUsecase(NewMockRepo())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetProductId := 1

	dataReq := entity.ProductRequest{
		Name:     "your step mom",
		Quantity: 1,
		Price:    100.0,
	}

	err := uc.UpdateProduct(ctx, targetProductId, &dataReq)
	if err != nil {
		t.Fatal(err)
	}

	products, err := uc.GetProducts(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if (*products)[0].Name != dataReq.Name || (*products)[0].Price != dataReq.Price || (*products)[0].Quantity != dataReq.Quantity {
		t.Errorf("expected %v, got %v", dataReq, (*products)[0])
	}
}

func TestDeleteProduct(t *testing.T) {
	uc := usecase.NewUsecase(NewMockRepo())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetProduct := mockProductData[0]

	err := uc.DeleteProduct(ctx, targetProduct.Id)
	if err != nil {
		t.Fatal(err)
	}

	products, err := uc.GetProducts(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if (*products)[0].Id == targetProduct.Id || (*products)[0].Name == targetProduct.Name {
		t.Errorf("expected delete but still exists, want delete %v", targetProduct)
	}
}
