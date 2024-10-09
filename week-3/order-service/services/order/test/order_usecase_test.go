package test

import (
	"context"
	"errors"
	"order_service/internal/core"
	orderEntity "order_service/services/order/entity"
	"order_service/services/order/test/mock"
	"order_service/services/order/usecase"
	productEntity "order_service/services/product/entity"
	userEntity "order_service/services/user/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type OrderUsecaseTestSuite struct {
	suite.Suite
	mockRepo *mock.MockOrderRepository
	usecase  usecase.OrderUsecase
}

func (suite *OrderUsecaseTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.mockRepo = mock.NewMockOrderRepository(ctrl)
	suite.usecase = usecase.NewUsecase(suite.mockRepo)
}

func (suite *OrderUsecaseTestSuite) TestCreateOrder() {
	items := []orderEntity.OrderItem{
		{
			OrderId:      1,
			ProductId:    1,
			ProductName:  "orange",
			ProductPrice: 25,
			Quantity:     2,
		},
		{
			OrderId:      1,
			ProductId:    2,
			ProductName:  "pineapple",
			ProductPrice: 50,
			Quantity:     1,
		},
	}

	tests := []struct {
		name      string
		order     *orderEntity.Order
		repoErr   error
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Successful order creation",
			order: &orderEntity.Order{
				Id:         1,
				UserId:     1,
				TotalPrice: 150,
				Items:      items,
				CreatedAt:  time.Now(),
			},
			repoErr:   nil,
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name: "Product out of stock",
			order: &orderEntity.Order{
				Id:         1,
				UserId:     1,
				TotalPrice: 150,
				Items:      items,
				CreatedAt:  time.Now(),
			},
			repoErr:   orderEntity.ErrOutOfStock,
			want:      core.ErrConfict.WithError(orderEntity.ErrOutOfStock.Error()),
			assertion: assert.Error,
		},
		{
			name: "Insufficient user's balance",
			order: &orderEntity.Order{
				Id:         1,
				UserId:     1,
				TotalPrice: 150,
				Items:      items,
				CreatedAt:  time.Now(),
			},
			repoErr:   orderEntity.ErrInsufficientBalance,
			want:      core.ErrConfict.WithError(orderEntity.ErrInsufficientBalance.Error()),
			assertion: assert.Error,
		},
		{
			name: "Unknown error",
			order: &orderEntity.Order{
				Id:         1,
				UserId:     1,
				TotalPrice: 150,
				Items:      items,
				CreatedAt:  time.Now(),
			},
			repoErr:   errors.New("this is an error"),
			want:      core.ErrInternalServerError.WithError(orderEntity.ErrCannotCreateOrder.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().CreateOrder(gomock.Any(), tt.order, gomock.Any()).Return(tt.repoErr)

			err := suite.usecase.CreateOrder(context.Background(), tt.order)

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.want, "error should return correctly")
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestCreateOrderCallback() {
	user := &userEntity.User{
		Id:       1,
		Username: "partridge",
		Password: "130703",
		Balance:  100,
	}

	products := &[]productEntity.Product{
		{
			Id:       1,
			Name:     "orange",
			Quantity: 2,
			Price:    25,
		},
		{
			Id:       2,
			Name:     "apple",
			Quantity: 1,
			Price:    50,
		},
	}

	tests := []struct {
		name      string
		order     *orderEntity.Order
		user      *userEntity.User
		products  *[]productEntity.Product
		want      bool
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Successful order validate",
			order: &orderEntity.Order{
				Id:         0,
				UserId:     1,
				TotalPrice: 100,
				Items: []orderEntity.OrderItem{
					{
						OrderId:      0,
						ProductId:    1,
						ProductName:  "orange",
						Quantity:     2,
						ProductPrice: 25,
					},
					{
						OrderId:      0,
						ProductId:    2,
						ProductName:  "apple",
						Quantity:     1,
						ProductPrice: 50,
					},
				},
			},
			user:      user,
			products:  products,
			want:      true,
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "Nil order",
			order:     nil,
			user:      user,
			products:  products,
			want:      false,
			wantErr:   orderEntity.ErrInvalidMemory,
			assertion: assert.Error,
		},
		{
			name: "Product's quantity is out of stock",
			order: &orderEntity.Order{
				Id:         0,
				UserId:     1,
				TotalPrice: 125,
				Items: []orderEntity.OrderItem{
					{
						OrderId:      0,
						ProductId:    1,
						ProductName:  "orange",
						Quantity:     3,
						ProductPrice: 25,
					},
					{
						OrderId:      0,
						ProductId:    2,
						ProductName:  "apple",
						Quantity:     1,
						ProductPrice: 50,
					},
				},
			},
			user:      user,
			products:  products,
			want:      false,
			wantErr:   orderEntity.ErrOutOfStock,
			assertion: assert.Error,
		},
		{
			name: "Order items length is not equal",
			order: &orderEntity.Order{
				Id:         0,
				UserId:     1,
				TotalPrice: 50,
				Items: []orderEntity.OrderItem{
					{
						OrderId:      0,
						ProductId:    1,
						ProductName:  "orange",
						Quantity:     2,
						ProductPrice: 25,
					},
				},
			},
			user:      user,
			products:  products,
			want:      false,
			wantErr:   orderEntity.ErrNotEqual,
			assertion: assert.Error,
		},
		{
			name: "Insufficient user's balance",
			order: &orderEntity.Order{
				Id:         0,
				UserId:     1,
				TotalPrice: 100,
				Items: []orderEntity.OrderItem{
					{
						OrderId:      0,
						ProductId:    1,
						ProductName:  "orange",
						Quantity:     2,
						ProductPrice: 25,
					},
					{
						OrderId:      0,
						ProductId:    2,
						ProductName:  "apple",
						Quantity:     1,
						ProductPrice: 50,
					},
				},
			},
			user: &userEntity.User{
				Id:       1,
				Username: "partridge",
				Password: "130703",
				Balance:  50,
			},
			products:  products,
			want:      false,
			wantErr:   orderEntity.ErrInsufficientBalance,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			accept, err := suite.usecase.CreateOrderCallback(tt.order, tt.user, tt.products)

			suite.Equal(tt.want, accept, "first return argument must be equal")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be return correctly")
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestGetOrders() {
	tests := []struct {
		name      string
		repoErr   error
		want      *[]orderEntity.Order
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:    "Successful get orders",
			repoErr: nil,
			want: &[]orderEntity.Order{
				{
					Id:         1,
					UserId:     1,
					TotalPrice: 100,
					Items: []orderEntity.OrderItem{
						{
							OrderId:      1,
							ProductId:    1,
							ProductName:  "orange",
							ProductPrice: 50,
							Quantity:     2,
						},
					},
				},
			},
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "Repo return an error",
			repoErr:   errors.New("this is an error"),
			want:      nil,
			wantErr:   core.ErrNotFound.WithError(orderEntity.ErrOrderNotFound.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().GetOrders(gomock.Any()).Return(tt.want, tt.repoErr)

			orders, err := suite.usecase.GetOrders(context.Background())

			suite.Equal(tt.want, orders, "orders should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be return correctly")
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestGetTopFiveOrdersByPrice() {
	tests := []struct {
		name      string
		repoErr   error
		want      *[]orderEntity.Order
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:    "Successful get orders",
			repoErr: nil,
			want: &[]orderEntity.Order{
				{
					Id:         1,
					UserId:     1,
					TotalPrice: 100,
					Items: []orderEntity.OrderItem{
						{
							OrderId:      1,
							ProductId:    1,
							ProductName:  "orange",
							ProductPrice: 50,
							Quantity:     2,
						},
					},
				},
			},
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "Order repo return an error",
			repoErr:   errors.New("this is an error"),
			want:      nil,
			wantErr:   core.ErrNotFound.WithError(orderEntity.ErrOrderNotFound.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().GetTopFiveOrdersByPrice(gomock.Any()).Return(tt.want, tt.repoErr)

			orders, err := suite.usecase.GetTopFiveOrdersByPrice(context.Background())

			suite.Equal(tt.want, orders, "orders should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be returned correctly")
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestGetNumOfOrdersByMonth() {
	tests := []struct {
		name      string
		userId    int
		repoErr   error
		want      *[]orderEntity.AggregatedOrdersByMonth
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:    "Successful get orders",
			userId:  1,
			repoErr: nil,
			want: &[]orderEntity.AggregatedOrdersByMonth{
				{
					NumOfOrders: 1,
					Time:        time.Now(),
				},
			},
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "Repo return an error",
			userId:    1,
			repoErr:   errors.New("this is an error"),
			want:      nil,
			wantErr:   core.ErrInternalServerError.WithError(orderEntity.ErrOrderNotFound.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().GetNumOfOrdersPerMonth(gomock.Any(), tt.userId).Return(tt.want, tt.repoErr)

			orders, err := suite.usecase.GetNumOfOrdersByMonth(context.Background(), tt.userId)

			suite.Equal(tt.want, orders, "orders should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be returned correctly")
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestGetOrdersSummarize() {
	tests := []struct {
		name               string
		startDate, endDate time.Time
		repoErr            error
		want               *[]orderEntity.OrdersSummarize
		wantErr            error
		assertion          assert.ErrorAssertionFunc
	}{
		{
			name:      "Successful get orders",
			startDate: time.Now(),
			endDate:   time.Now().Add(24 * time.Hour),
			repoErr:   nil,
			want: &[]orderEntity.OrdersSummarize{
				{
					UserId:                   1,
					Username:                 "partridge",
					NumOfOrders:              1,
					SumOrderPrice:            25,
					AverageOrderItemQuantity: 2,
				},
			},
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "Repo return an error",
			startDate: time.Now(),
			endDate:   time.Now().Add(24 * time.Hour),
			repoErr:   errors.New("this is an error"),
			want:      nil,
			wantErr:   core.ErrInternalServerError.WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().GetOrdersSummarize(gomock.Any(), tt.startDate, tt.endDate).Return(tt.want, tt.repoErr)

			orders, err := suite.usecase.GetOrdersSummarize(context.Background(), tt.startDate, tt.endDate)

			suite.Equal(tt.want, orders, "orders should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be returned correctly")
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestGetOrder() {
	tests := []struct {
		name            string
		userId, orderId int
		repoErr         error
		want            *orderEntity.Order
		wantErr         error
		assertion       assert.ErrorAssertionFunc
	}{
		{
			name:    "Successful get order",
			userId:  1,
			orderId: 1,
			repoErr: nil,
			want: &orderEntity.Order{
				Id:         1,
				UserId:     1,
				TotalPrice: 100,
				Items: []orderEntity.OrderItem{
					{
						OrderId:      1,
						ProductId:    1,
						ProductName:  "orange",
						ProductPrice: 50,
						Quantity:     2,
					},
				},
			},
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "Repo return an error",
			userId:    1,
			orderId:   1,
			repoErr:   errors.New("this is an error"),
			want:      nil,
			wantErr:   core.ErrNotFound.WithError(orderEntity.ErrOrderNotFound.Error()).WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().GetOrder(gomock.Any(), tt.userId, tt.orderId).Return(tt.want, tt.repoErr)

			orders, err := suite.usecase.GetOrder(context.Background(), tt.userId, tt.orderId)

			suite.Equal(tt.want, orders, "order should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be returned correctly")
			}
		})
	}
}

func TestOrderUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderUsecaseTestSuite))
}
