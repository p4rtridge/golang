package test

import (
	"order_service/services/order/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OrderTestSuite struct {
	suite.Suite
	order entity.Order
}

func (suite *OrderTestSuite) SetupTest() {
	suite.order = entity.Order{
		Id:         1,
		UserId:     1,
		TotalPrice: 100,
		Items: []entity.OrderItem{
			{
				OrderId:      1,
				ProductId:    1,
				ProductName:  "orange",
				Quantity:     1,
				ProductPrice: 50,
			},
			{
				OrderId:      1,
				ProductId:    2,
				ProductName:  "lemon",
				Quantity:     2,
				ProductPrice: 25,
			},
		},
	}
}

func (suite *OrderTestSuite) TestNewOrder() {
	items := []entity.OrderItem{
		{
			OrderId:      1,
			ProductId:    1,
			ProductName:  "orange",
			Quantity:     1,
			ProductPrice: 50,
		},
		{
			OrderId:      1,
			ProductId:    2,
			ProductName:  "lemon",
			Quantity:     2,
			ProductPrice: 25,
		},
	}

	order := entity.NewOrder(1, 1, 100, items)

	assert.Equal(suite.T(), 1, order.Id, "Id should be set correctly")
	assert.Equal(suite.T(), 1, order.UserId, "User Id should be set correctly")
	assert.Equal(suite.T(), float32(100), order.TotalPrice, "TotalPrice should be set correctly")
	assert.Equal(suite.T(), items, order.Items, "Order's Items should be set correctly")
}

func (suite *OrderTestSuite) TestSetId() {
	tests := []struct {
		name  string
		order *entity.Order
		id    int
		want  int
		panic bool
	}{
		{
			name:  "Non nil order",
			order: &suite.order,
			id:    2,
			want:  2,
			panic: false,
		},
		{
			name:  "Nil order",
			order: nil,
			id:    2,
			want:  0,
			panic: true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			assert.NotPanics(suite.T(), func() {
				tt.order.SetId(tt.id)
			}, "Calling SetId on nil order should not be panic")
		} else {
			tt.order.SetId(tt.id)

			assert.Equal(suite.T(), tt.want, tt.order.Id, "Id should be updated")
		}
	}
}

func (suite *OrderTestSuite) TestSetUserId() {
	tests := []struct {
		name   string
		order  *entity.Order
		userId int
		want   int
		panic  bool
	}{
		{
			name:   "Non nil order",
			order:  &suite.order,
			userId: 2,
			want:   2,
			panic:  false,
		},
		{
			name:   "Nil order",
			order:  nil,
			userId: 2,
			want:   0,
			panic:  true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			assert.NotPanics(suite.T(), func() {
				tt.order.SetUserId(tt.userId)
			}, "Calling SetUserId on nil order should not be panic")
		} else {
			tt.order.SetUserId(tt.userId)

			assert.Equal(suite.T(), tt.want, tt.order.UserId, "UserId should be updated")
		}
	}
}

func (suite *OrderTestSuite) TestSetTotalPrice() {
	tests := []struct {
		name       string
		order      *entity.Order
		totalPrice float32
		want       float32
		panic      bool
	}{
		{
			name:       "Non nil order",
			order:      &suite.order,
			totalPrice: 200,
			want:       200,
			panic:      false,
		},
		{
			name:       "Nil order",
			order:      nil,
			totalPrice: 200,
			want:       0,
			panic:      true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			assert.NotPanics(suite.T(), func() {
				tt.order.SetTotalPrice(tt.totalPrice)
			}, "Calling SetTotalPrice on nil order should not be panic")
		} else {
			tt.order.SetTotalPrice(tt.totalPrice)

			assert.Equal(suite.T(), tt.want, tt.order.TotalPrice, "TotalPrice should be updated")
		}
	}
}

func (suite *OrderTestSuite) TestSetCreatedAt() {
	now := time.Now()

	tests := []struct {
		name      string
		order     *entity.Order
		createdAt time.Time
		want      time.Time
		panic     bool
	}{
		{
			name:      "Non nil order",
			order:     &suite.order,
			createdAt: now,
			want:      now,
			panic:     false,
		},
		{
			name:      "Nil order",
			order:     nil,
			createdAt: now,
			want:      time.Time{},
			panic:     true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			assert.NotPanics(suite.T(), func() {
				tt.order.SetCreatedAt(tt.createdAt)
			}, "Calling SetCreatedAt on nil order should not be panic")
		} else {
			tt.order.SetCreatedAt(tt.createdAt)

			assert.Equal(suite.T(), tt.want, tt.order.CreatedAt, "CreatedAt should be updated")
		}
	}
}

func (suite *OrderTestSuite) TestSetUpdatedAt() {
	now := time.Now()

	tests := []struct {
		name      string
		order     *entity.Order
		updatedAt time.Time
		want      time.Time
		panic     bool
	}{
		{
			name:      "Non nil order",
			order:     &suite.order,
			updatedAt: now,
			want:      now,
			panic:     false,
		},
		{
			name:      "Nil order",
			order:     nil,
			updatedAt: now,
			want:      time.Time{},
			panic:     true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			assert.NotPanics(suite.T(), func() {
				tt.order.SetUpdatedAt(&tt.updatedAt)
			}, "Calling SetUpdatedAt on nil order should not be panic")
		} else {
			tt.order.SetUpdatedAt(&tt.updatedAt)

			assert.Equal(suite.T(), tt.want, *tt.order.UpdatedAt, "UpdatedAt should be updated")
		}
	}
}

func (suite *OrderTestSuite) TestAddItem_NonNilOrder() {
	item := entity.OrderItem{
		OrderId:      1,
		ProductId:    3,
		ProductName:  "pineapple",
		ProductPrice: 25.0,
		Quantity:     2,
	}

	tests := []struct {
		name  string
		order *entity.Order
		item  entity.OrderItem
		want  entity.OrderItem
		panic bool
	}{
		{
			name:  "Non nil order",
			order: &suite.order,
			item:  item,
			want:  item,
			panic: false,
		},
		{
			name:  "Nil order",
			order: nil,
			item:  item,
			want:  entity.OrderItem{},
			panic: true,
		},
	}

	for _, tt := range tests {
		if tt.panic {
			assert.NotPanics(suite.T(), func() {
				tt.order.AddItem(tt.item)
			}, "Calling AddItem on nil order should not be panic")
		} else {
			tt.order.AddItem(tt.item)

			assert.Equal(suite.T(), tt.want, tt.order.Items[len(tt.order.Items)-1], "Item should be updated")
		}
	}
}

func (suite *OrderTestSuite) TestGetIdSafe() {
	tests := []struct {
		name  string
		order *entity.Order
		want  int
	}{
		{
			name:  "Non nil order",
			order: &suite.order,
			want:  1,
		},
		{
			name:  "Nil order",
			order: nil,
			want:  0,
		},
	}

	for _, tt := range tests {
		got := tt.order.GetIdSafe()

		assert.Equal(suite.T(), tt.want, got, "Id should be retrieved correctly")
	}
}

func (suite *OrderTestSuite) TestGetUserIdSafe() {
	tests := []struct {
		name  string
		order *entity.Order
		want  int
	}{
		{
			name:  "Non nil order",
			order: &suite.order,
			want:  1,
		},
		{
			name:  "Nil order",
			order: nil,
			want:  0,
		},
	}

	for _, tt := range tests {
		got := tt.order.GetUserIdSafe()

		assert.Equal(suite.T(), tt.want, got, "UserId should be retrieved correctly")
	}
}

func (suite *OrderTestSuite) TestGetTotalPriceSafe() {
	tests := []struct {
		name  string
		order *entity.Order
		want  float32
	}{
		{
			name:  "Non nil order",
			order: &suite.order,
			want:  100,
		},
		{
			name:  "Nil order",
			order: nil,
			want:  0,
		},
	}

	for _, tt := range tests {
		got := tt.order.GetTotalPriceSafe()

		assert.Equal(suite.T(), tt.want, got, "TotalPrice should be retrieved correctly")
	}
}

func (suite *OrderTestSuite) TestGetItemsSafe() {
	tests := []struct {
		name  string
		order *entity.Order
		want  []entity.OrderItem
	}{
		{
			name:  "Non nil order",
			order: &suite.order,
			want:  suite.order.Items,
		},
		{
			name:  "Nil order",
			order: nil,
			want:  []entity.OrderItem{},
		},
	}

	for _, tt := range tests {
		got := tt.order.GetItemsSafe()

		assert.Equal(suite.T(), tt.want, got, "Items should be retrieved correctly")
	}
}

func (suite *OrderTestSuite) TestGetItemSafe() {
	tests := []struct {
		name  string
		order *entity.Order
		index int
		want  *entity.OrderItem
	}{
		{
			name:  "Non nil order",
			order: &suite.order,
			index: 0,
			want:  &suite.order.Items[0],
		},
		{
			name:  "Nil order",
			order: nil,
			index: 0,
			want:  nil,
		},
	}

	for _, tt := range tests {
		got := tt.order.GetItemSafe(tt.index)

		assert.Equal(suite.T(), tt.want, got, "Item should be retrieved correctly")
	}
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}
