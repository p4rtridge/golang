package entity_test

import (
	"order_service/services/product/entity"
	"testing"
)

func TestNewProduct(t *testing.T) {
	mock_product := entity.Product{
		Id:       1,
		Name:     "your mom",
		Quantity: 1,
		Price:    100.0,
	}

	product := entity.NewProduct(1, "your mom", 1, 100.0)

	if product.Id != mock_product.Id || product.Name != mock_product.Name || product.Quantity != mock_product.Quantity || product.Price != mock_product.Price {
		t.Errorf("expected %v, got %v", mock_product, product)
	}
}

func TestSetId(t *testing.T) {
	// case non nil
	expectedId := 1

	var product *entity.Product = &entity.Product{
		Id:       0,
		Name:     "your mom",
		Quantity: 1,
		Price:    100.0,
	}

	product.SetId(expectedId)

	if product.Id != expectedId {
		t.Errorf("expected %v, got %v", expectedId, product.Id)
	}

	// case nil
	product = nil

	product.SetId(expectedId)

	if product != nil && product.Id == expectedId {
		t.Errorf("expected %v, got %v", nil, product.Id)
	}
}

func TestSetQuantity(t *testing.T) {
	// case non nil
	expectedQuantity := 10

	var product *entity.Product = &entity.Product{
		Id:       1,
		Name:     "your mom",
		Quantity: 1,
		Price:    100.0,
	}

	product.SetQuantity(expectedQuantity)

	if product.Quantity != expectedQuantity {
		t.Errorf("expected %v, got %v", expectedQuantity, product.Quantity)
	}

	// case nil
	product = nil

	product.SetQuantity(expectedQuantity)

	if product != nil && product.Quantity == expectedQuantity {
		t.Errorf("expected %v, got %v", nil, product.Quantity)
	}
}

func TestGetName(t *testing.T) {
	expectedName := "your mom"

	product := entity.Product{
		Id:       1,
		Name:     expectedName,
		Quantity: 1,
		Price:    100.0,
	}

	if product.GetName() != expectedName {
		t.Errorf("expected %v, got %v", expectedName, product.GetName())
	}
}

func TestGetQuantity(t *testing.T) {
	expectedQuantity := 10

	product := entity.Product{
		Id:       1,
		Name:     "your mom",
		Quantity: expectedQuantity,
		Price:    100.0,
	}

	if product.GetQuantity() != expectedQuantity {
		t.Errorf("expected %v, got %v", expectedQuantity, product.GetQuantity())
	}
}

func TestGetPrice(t *testing.T) {
	expectedPrice := float32(100.0)

	product := entity.Product{
		Id:       1,
		Name:     "your mom",
		Quantity: 1,
		Price:    expectedPrice,
	}

	if product.GetPrice() != expectedPrice {
		t.Errorf("expected %v, got %v", expectedPrice, product.GetPrice())
	}
}
