package entity_test

import (
	"order_service/services/order/entity"
	"testing"
)

func TestNewOrderItem(t *testing.T) {
	mock_order_item := entity.OrderItem{
		OrderId:      1,
		ProductId:    1,
		ProductName:  "your mom",
		ProductPrice: 100.0,
		Quantity:     1,
	}

	orderItem := entity.NewOrderItem(1, 1, "your mom", 100.0, 1)

	if orderItem != mock_order_item {
		t.Errorf("expected %v, got %v", mock_order_item, orderItem)
	}
}

func TestSetOrderId(t *testing.T) {
	// case non nil
	expectedOrderId := 1

	var orderItem *entity.OrderItem = &entity.OrderItem{
		OrderId:      0,
		ProductId:    1,
		ProductName:  "your mom",
		ProductPrice: 100.0,
		Quantity:     1,
	}

	orderItem.SetOrderId(expectedOrderId)

	if orderItem.OrderId != expectedOrderId {
		t.Errorf("expected %v, got %v", expectedOrderId, orderItem.OrderId)
	}

	// case nil
	orderItem = nil

	orderItem.SetOrderId(expectedOrderId)

	if orderItem != nil && orderItem.OrderId == expectedOrderId {
		t.Errorf("expected %v, got %v", nil, orderItem.OrderId)
	}
}

func TestSetProductId(t *testing.T) {
	// case non nil
	expectedProductId := 1

	var orderItem *entity.OrderItem = &entity.OrderItem{
		OrderId:      1,
		ProductId:    0,
		ProductName:  "your mom",
		ProductPrice: 100.0,
		Quantity:     1,
	}

	orderItem.SetProductId(expectedProductId)

	if orderItem.ProductId != expectedProductId {
		t.Errorf("expected %v, got %v", expectedProductId, orderItem.ProductId)
	}

	// case nil
	orderItem = nil

	orderItem.SetOrderId(expectedProductId)

	if orderItem != nil && orderItem.ProductId == expectedProductId {
		t.Errorf("expected %v, got %v", nil, orderItem.ProductId)
	}
}

func TestSetProductName(t *testing.T) {
	// case non nil
	expectedProductName := "your mom"

	var orderItem *entity.OrderItem = &entity.OrderItem{
		OrderId:      1,
		ProductId:    1,
		ProductName:  "your dad",
		ProductPrice: 100.0,
		Quantity:     1,
	}

	orderItem.SetProductName(expectedProductName)

	if orderItem.ProductName != expectedProductName {
		t.Errorf("expected %v, got %v", expectedProductName, orderItem.ProductName)
	}

	// case nil
	orderItem = nil

	orderItem.SetProductName(expectedProductName)

	if orderItem != nil && orderItem.ProductName == expectedProductName {
		t.Errorf("expected %v, got %v", nil, orderItem.ProductName)
	}
}

func TestSetProductPrice(t *testing.T) {
	// case non nil
	expectedProductPrice := float32(100.0)

	var orderItem *entity.OrderItem = &entity.OrderItem{
		OrderId:      1,
		ProductId:    1,
		ProductName:  "your dad",
		ProductPrice: 10.0,
		Quantity:     1,
	}

	orderItem.SetProductPrice(expectedProductPrice)

	if orderItem.ProductPrice != expectedProductPrice {
		t.Errorf("expected %v, got %v", expectedProductPrice, orderItem.ProductPrice)
	}

	// case nil
	orderItem = nil

	orderItem.SetProductPrice(expectedProductPrice)

	if orderItem != nil && orderItem.ProductPrice == expectedProductPrice {
		t.Errorf("expected %v, got %v", nil, orderItem.ProductPrice)
	}
}

func TestGetProductId(t *testing.T) {
	expectedProductId := 1

	orderItem := entity.OrderItem{
		OrderId:      1,
		ProductId:    expectedProductId,
		ProductName:  "your mom",
		ProductPrice: 100.0,
		Quantity:     1,
	}

	if orderItem.GetProductId() != expectedProductId {
		t.Errorf("expected %v, got %v", expectedProductId, orderItem.GetProductId())
	}
}

func TestGetQuantity(t *testing.T) {
	expectedQuantity := 1

	orderItem := entity.OrderItem{
		OrderId:      1,
		ProductId:    1,
		ProductName:  "your mom",
		ProductPrice: 100.0,
		Quantity:     expectedQuantity,
	}

	if orderItem.GetQuantity() != expectedQuantity {
		t.Errorf("expected %v, got %v", expectedQuantity, orderItem.GetQuantity())
	}
}

func TestGetProductName(t *testing.T) {
	expectedProductName := "your mom"

	orderItem := entity.OrderItem{
		OrderId:      1,
		ProductId:    1,
		ProductName:  expectedProductName,
		ProductPrice: 100.0,
		Quantity:     1,
	}

	if orderItem.GetProductName() != expectedProductName {
		t.Errorf("expected %v, got %v", expectedProductName, orderItem.GetProductName())
	}
}

func TestGetProductPrice(t *testing.T) {
	expectedProductPrice := float32(100.0)

	orderItem := entity.OrderItem{
		OrderId:      1,
		ProductId:    1,
		ProductName:  "your mom",
		ProductPrice: expectedProductPrice,
		Quantity:     1,
	}

	if orderItem.GetProductPrice() != expectedProductPrice {
		t.Errorf("expected %v, got %v", expectedProductPrice, orderItem.GetProductPrice())
	}
}
