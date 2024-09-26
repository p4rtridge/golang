package test

import (
	"order_service/services/order/entity"
	"testing"

	"github.com/stretchr/testify/suite"
)

type OrderItemTestSuite struct {
	suite.Suite
	orderItem entity.OrderItem
}

func (suite *OrderItemTestSuite) SetupTest() {
	suite.orderItem = entity.OrderItem{
		OrderId:      1,
		ProductId:    1,
		ProductName:  "orange",
		Quantity:     1,
		ProductPrice: 50,
	}
}

func (suite *OrderItemTestSuite) TestNewOrderItem() {
	item := entity.NewOrderItem(1, 1, "orange", 50, 1)

	suite.Equal(1, item.OrderId, "OrderId should be set correctly")
	suite.Equal(1, item.ProductId, "ProductId should be set correctly")
	suite.Equal("orange", item.ProductName, "ProductName should be set correctly")
	suite.Equal(float32(50), item.ProductPrice, "ProductPrice should be set correctly")
	suite.Equal(1, item.Quantity, "Quantity should be set correctly")
}

func (suite *OrderItemTestSuite) TestSetOrderId() {
	tests := []struct {
		name      string
		orderItem *entity.OrderItem
		orderId   int
		want      int
		panic     bool
	}{
		{
			name:      "Non nil order item",
			orderItem: &suite.orderItem,
			orderId:   2,
			want:      2,
			panic:     false,
		},
		{
			name:      "Nil order item",
			orderItem: nil,
			orderId:   2,
			want:      0,
			panic:     true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			suite.NotPanics(func() {
				tt.orderItem.SetOrderId(tt.orderId)
			}, "Calling SetOrderId on nil order item should not be panic")
		} else {
			tt.orderItem.SetOrderId(tt.orderId)

			suite.Equal(tt.want, suite.orderItem.OrderId, "OrderId should be updated")
		}
	}
}

func (suite *OrderItemTestSuite) TestSetProductId() {
	tests := []struct {
		name      string
		orderItem *entity.OrderItem
		productId int
		want      int
		panic     bool
	}{
		{
			name:      "Non nil order item",
			orderItem: &suite.orderItem,
			productId: 2,
			want:      2,
			panic:     false,
		},
		{
			name:      "Nil order item",
			orderItem: nil,
			productId: 2,
			want:      0,
			panic:     true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			suite.NotPanics(func() {
				tt.orderItem.SetProductId(tt.productId)
			}, "Calling SetProductId on nil order item should not be panic")
		} else {
			tt.orderItem.SetProductId(tt.productId)

			suite.Equal(tt.want, tt.orderItem.ProductId, "ProductId should be updated")
		}
	}
}

func (suite *OrderItemTestSuite) TestSetProductName() {
	tests := []struct {
		name        string
		orderItem   *entity.OrderItem
		productName string
		want        string
		panic       bool
	}{
		{
			name:        "Non nil order item",
			orderItem:   &suite.orderItem,
			productName: "pineapple",
			want:        "pineapple",
			panic:       false,
		},
		{
			name:        "Nil order item",
			orderItem:   nil,
			productName: "pineapple",
			want:        "",
			panic:       true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			suite.NotPanics(func() {
				tt.orderItem.SetProductName(tt.productName)
			}, "Calling SetProductName on nil order item should not be panic")
		} else {
			tt.orderItem.SetProductName(tt.productName)

			suite.Equal(tt.want, tt.orderItem.ProductName, "ProductName should be updated")
		}
	}
}

func (suite *OrderItemTestSuite) TestSetProductPrice() {
	tests := []struct {
		name         string
		orderItem    *entity.OrderItem
		productPrice float32
		want         float32
		panic        bool
	}{
		{
			name:         "Non nil order item",
			orderItem:    &suite.orderItem,
			productPrice: 100,
			want:         100,
			panic:        false,
		},
		{
			name:         "Nil order item",
			orderItem:    nil,
			productPrice: 100,
			want:         0,
			panic:        true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			suite.NotPanics(func() {
				tt.orderItem.SetProductPrice(tt.productPrice)
			}, "Calling SetProductPrice on nil order item should not be panic")
		} else {
			tt.orderItem.SetProductPrice(tt.productPrice)

			suite.Equal(tt.want, tt.orderItem.ProductPrice, "ProductPrice should be updated")
		}
	}
}

func (suite *OrderItemTestSuite) TestGetProductId() {
	got := suite.orderItem.GetProductId()

	suite.Equal(suite.orderItem.ProductId, got, "ProductPrice should be retrieved correctly")
}

func (suite *OrderItemTestSuite) TestGetQuantity() {
	got := suite.orderItem.GetQuantity()

	suite.Equal(suite.orderItem.Quantity, got, "Quantity should be retrieved correctly")
}

func (suite *OrderItemTestSuite) TestProductName() {
	got := suite.orderItem.GetProductName()

	suite.Equal(suite.orderItem.ProductName, got, "ProductName should be retrieved correctly")
}

func (suite *OrderItemTestSuite) TestProductPrice() {
	got := suite.orderItem.GetProductPrice()

	suite.Equal(suite.orderItem.ProductPrice, got, "ProductPrice should be retrieved correctly")
}

func TestOrderItemTestSuite(t *testing.T) {
	suite.Run(t, new(OrderItemTestSuite))
}
