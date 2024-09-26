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
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := tt.order.Validate()

			if tt.assertion(suite.T(), err) {
				assert.ErrorIs(suite.T(), err, tt.want, "error should be return correctly")
			}
		})
	}
}

func TestOrderVarsTestSuite(t *testing.T) {
	suite.Run(t, new(OrderVarsTestSuite))
}
