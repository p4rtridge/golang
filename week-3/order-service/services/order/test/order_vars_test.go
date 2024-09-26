package test

import (
	"order_service/services/order/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OrderVarsTestSuite struct {
	suite.Suite
}

func (suite *OrderVarsTestSuite) TestValidate() {
	tests := []struct {
		name      string
		order     entity.OrderRequest
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Validated order request",
			order: entity.OrderRequest{
				Items: []entity.ProductItem{
					{
						ProductId: 1,
						Quantity:  1,
					},
				},
			},
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name: "Empty order items",
			order: entity.OrderRequest{
				Items: nil,
			},
			want:      entity.ErrItemEmpty,
			assertion: assert.Error,
		},
		{
			name: "Zero value order item",
			order: entity.OrderRequest{
				Items: []entity.ProductItem{
					{
						ProductId: 0,
						Quantity:  1,
					},
				},
			},
			want:      entity.ErrMissingField,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := tt.order.Validate()

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.want, "error should be return correctly")
			}
		})
	}
}

func (suite *OrderVarsTestSuite) TestGetItems() {
	orderReq := entity.OrderRequest{
		Items: []entity.ProductItem{
			{
				ProductId: 1,
				Quantity:  1,
			},
		},
	}

	tests := []struct {
		name  string
		order entity.OrderRequest
		want  []entity.ProductItem
	}{
		{
			name:  "Order items exists",
			order: orderReq,
			want:  orderReq.Items,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			items := tt.order.GetItems()

			suite.Equal(tt.want, items, "items should be retrieved correctly")
		})
	}
}

func (suite *OrderVarsTestSuite) TestGetItemId() {
	tests := []struct {
		name string
		item entity.ProductItem
		want int
	}{
		{
			name: "Order items exists",
			item: entity.ProductItem{
				ProductId: 1,
				Quantity:  1,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			id := tt.item.GetItemId()

			suite.Equal(tt.want, id, "id should be retrieved correctly")
		})
	}
}

func (suite *OrderVarsTestSuite) TestGetItemQuantity() {
	tests := []struct {
		name string
		item entity.ProductItem
		want int
	}{
		{
			name: "Order items exists",
			item: entity.ProductItem{
				ProductId: 1,
				Quantity:  1,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			quantity := tt.item.GetItemQuantity()

			suite.Equal(tt.want, quantity, "quantity should be retrieved correctly")
		})
	}
}

func TestOrderVarsTestSuite(t *testing.T) {
	suite.Run(t, new(OrderVarsTestSuite))
}
