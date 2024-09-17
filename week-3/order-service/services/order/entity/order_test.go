package entity_test

import (
	"order_service/services/order/entity"
	"testing"
)

var items []entity.OrderItem = []entity.OrderItem{
	{
		OrderId:      1,
		ProductId:    1,
		ProductName:  "your mom",
		ProductPrice: 100.0,
		Quantity:     1,
	},
}

func TestNewOrder(t *testing.T) {
	mock_order := entity.Order{
		Id:         1,
		UserId:     1,
		TotalPrice: 100.0,
		Items:      items,
	}

	order := entity.NewOrder(1, 1, 100.0, items)

	if order.Id != mock_order.Id || order.UserId != mock_order.UserId || order.TotalPrice != mock_order.TotalPrice {
		t.Errorf("expected %v, got %v", mock_order, order)
	}

	for idx, item := range order.Items {
		if item != mock_order.Items[idx] {
			t.Errorf("expected %v, got %v", mock_order, order)
		}
	}
}

func TestSetId(t *testing.T) {
	expectedId := 1

	order := entity.Order{
		Id:         0,
		UserId:     1,
		TotalPrice: 100.0,
		Items:      items,
	}

	order.SetId(expectedId)

	if order.Id != expectedId {
		t.Errorf("expected %v, got %v", expectedId, order.Id)
	}
}

func TestSetTotalPrice(t *testing.T) {
	expectedTP := float32(100.0)

	order := entity.Order{
		Id:         1,
		UserId:     1,
		TotalPrice: 10.0,
		Items:      items,
	}

	order.SetTotalPrice(expectedTP)

	if order.TotalPrice != expectedTP {
		t.Errorf("expected %v, got %v", expectedTP, order.Id)
	}
}

func TestGetIdSafe(t *testing.T) {
	// non nil
	expectedId := 1

	var order *entity.Order = &entity.Order{
		Id:         expectedId,
		UserId:     1,
		TotalPrice: 10.0,
		Items:      items,
	}

	if order.GetIdSafe() != expectedId {
		t.Errorf("expected %v, got %v", expectedId, order.GetIdSafe())
	}

	// nil
	expectedId = 0

	order = nil

	if order.GetIdSafe() != expectedId {
		t.Errorf("expected %v, got %v", expectedId, order.GetIdSafe())
	}
}

func TestGetUserIdSafe(t *testing.T) {
	// non nil
	expectedUserId := 1

	var order *entity.Order = &entity.Order{
		Id:         1,
		UserId:     expectedUserId,
		TotalPrice: 10.0,
		Items:      items,
	}

	if order.GetUserIdSafe() != expectedUserId {
		t.Errorf("expected %v, got %v", expectedUserId, order.GetUserIdSafe())
	}

	// nil
	expectedUserId = 0

	order = nil

	if order.GetIdSafe() != expectedUserId {
		t.Errorf("expected %v, got %v", expectedUserId, order.GetUserIdSafe())
	}
}

func TestGetTotalPriceSafe(t *testing.T) {
	// non nil
	expectedTotalPrice := float32(100.0)

	var order *entity.Order = &entity.Order{
		Id:         1,
		UserId:     1,
		TotalPrice: expectedTotalPrice,
		Items:      items,
	}

	if order.GetTotalPriceSafe() != expectedTotalPrice {
		t.Errorf("expected %v, got %v", expectedTotalPrice, order.GetTotalPriceSafe())
	}

	// nil
	expectedTotalPrice = 0

	order = nil

	if order.GetTotalPriceSafe() != expectedTotalPrice {
		t.Errorf("expected %v, got %v", expectedTotalPrice, order.GetTotalPriceSafe())
	}
}

func TestGetItemsSafe(t *testing.T) {
	// non nil
	expectedItems := []entity.OrderItem{
		{
			OrderId:      1,
			ProductId:    1,
			ProductName:  "your mom",
			ProductPrice: 100.0,
			Quantity:     1,
		},
	}

	var order *entity.Order = &entity.Order{
		Id:         1,
		UserId:     1,
		TotalPrice: 100.0,
		Items:      expectedItems,
	}

	if len(order.GetItemsSafe()) != len(expectedItems) {
		t.Errorf("expected %v items, got %v items", len(order.GetItemsSafe()), len(expectedItems))
	}

	for idx, item := range order.GetItemsSafe() {
		if expectedItems[idx] != item {
			t.Errorf("expected item %v in order, got %v", expectedItems[idx], item)
		}
	}

	// nil
	expectedItems = []entity.OrderItem{}

	order = nil

	if len(order.GetItemsSafe()) != len(expectedItems) {
		t.Errorf("expected %v items, got %v items", len(expectedItems), len(order.GetItemsSafe()))
	}
}

func TestGetItemSafe(t *testing.T) {
	// non nil
	expectedItems := []entity.OrderItem{
		{
			OrderId:      1,
			ProductId:    1,
			ProductName:  "your mom",
			ProductPrice: 100.0,
			Quantity:     1,
		},
	}

	var order *entity.Order = &entity.Order{
		Id:         1,
		UserId:     1,
		TotalPrice: 100.0,
		Items:      expectedItems,
	}

	if order.GetItemSafe(-1) != nil {
		t.Errorf("expected %v, got %v", nil, order.GetItemSafe(-1))
	}

	if order.GetItemSafe(0) == nil {
		t.Errorf("expected %v, got %v", expectedItems[0], order.GetItemSafe(0))
	}

	if order.GetItemSafe(1) != nil {
		t.Errorf("expected %v, got %v", nil, order.GetItemSafe(1))
	}

	// nil
	expectedItems = nil

	order = nil

	if order.GetItemSafe(0) != nil {
		t.Errorf("expected %v, got %v", nil, order.GetItemSafe(0))
	}
}
