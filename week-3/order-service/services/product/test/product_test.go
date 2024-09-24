package test

import (
	"order_service/services/product/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProductTestSuite struct {
	suite.Suite

	product entity.Product
}

func (suite *ProductTestSuite) SetupTest() {
	suite.product = entity.Product{
		Id:        1,
		Name:      "orange",
		Quantity:  10,
		Price:     2.5,
		CreatedAt: time.Now(),
	}
}

func (suite *ProductTestSuite) TestNewProduct() {
	product := entity.NewProduct(1, "orange", 10, 2.5)

	assert.Equal(suite.T(), 1, product.Id, "Id should be set correctly")
	assert.Equal(suite.T(), "orange", product.Name, "Name should be set correctly")
	assert.Equal(suite.T(), 10, product.Quantity, "Quantity should be set correctly")
	assert.Equal(suite.T(), float32(2.5), product.Price, "Price should be set correctly")
}

func (suite *ProductTestSuite) TestSetId() {
	tests := []struct {
		name    string
		product *entity.Product
		id      int
		want    int
		panic   bool
	}{
		{
			name:    "Non nil product",
			product: &suite.product,
			id:      2,
			want:    2,
			panic:   false,
		},
		{
			name:    "Nil product",
			product: nil,
			id:      2,
			want:    0,
			panic:   true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			suite.Run(tt.name, func() {
				assert.NotPanics(suite.T(), func() {
					tt.product.SetId(tt.id)
				}, "Calling SetId on nil product should not panic")
			})
		} else {
			suite.Run(tt.name, func() {
				tt.product.SetId(tt.id)

				assert.Equal(suite.T(), tt.want, tt.product.Id, "Id should be updated")
			})
		}
	}
}

func (suite *ProductTestSuite) TestSetQuantity() {
	tests := []struct {
		name     string
		product  *entity.Product
		quantity int
		want     int
		panic    bool
	}{
		{
			name:     "Non nil product",
			product:  &suite.product,
			quantity: 5,
			want:     5,
			panic:    false,
		},
		{
			name:     "Nil product",
			product:  nil,
			quantity: 5,
			want:     5,
			panic:    true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			suite.Run(tt.name, func() {
				assert.NotPanics(suite.T(), func() {
					tt.product.SetQuantity(tt.quantity)
				}, "Calling SetQuantity on nil should no panic")
			})
		} else {
			tt.product.SetQuantity(tt.quantity)

			assert.Equal(suite.T(), tt.want, tt.product.Quantity, "Quantity should be updated")
		}
	}
}

func (suite *ProductTestSuite) TestGetId() {
	got := suite.product.GetId()

	assert.Equal(suite.T(), suite.product.Id, got, "Id should be retrieved correctly")
}

func (suite *ProductTestSuite) TestGetName() {
	got := suite.product.GetName()

	assert.Equal(suite.T(), suite.product.Name, got, "Name should be retrieved correctly")
}

func (suite *ProductTestSuite) TestGetQuantity() {
	got := suite.product.GetQuantity()

	assert.Equal(suite.T(), suite.product.Quantity, got, "Quantity should be retrieved correctly")
}

func (suite *ProductTestSuite) TestGetPrice() {
	got := suite.product.GetPrice()

	assert.Equal(suite.T(), suite.product.Price, got, "Price should be retrieved correctly")
}

func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}
