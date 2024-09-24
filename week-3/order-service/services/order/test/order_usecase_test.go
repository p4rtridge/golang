package test

import (
	"context"
	orderEntity "order_service/services/order/entity"
	productEntity "order_service/services/product/entity"
	userEntity "order_service/services/user/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockOrderRepo struct {
	mock.Mock
}

func (mock *mockOrderRepo) CreateOrder(ctx context.Context, order *orderEntity.Order, callbackFn func(order *orderEntity.Order, user *userEntity.User, products productEntity.Product) (bool, error)) error {
	args := mock.Called(ctx, order, callbackFn)

	return args.Error(0)
}

func (mock *mockOrderRepo) GetOrders(ctx context.Context, userId int) (*[]orderEntity.Order, error) {
	args := mock.Called(ctx, userId)

	return args.Get(0).(*[]orderEntity.Order), args.Error(1)
}

func (mock *mockOrderRepo) GetOrdersSummarize(ctx context.Context, startDate, endDate time.Time) (*[]orderEntity.OrdersSummarize, error) {
	args := mock.Called(ctx, startDate, endDate)

	return args.Get(0).(*[]orderEntity.OrdersSummarize), args.Error(1)
}

func (mock *mockOrderRepo) GetTopFiveOrdersByPrice(ctx context.Context) (*[]orderEntity.Order, error) {
	args := mock.Called(ctx)

	return args.Get(0).(*[]orderEntity.Order), args.Error(1)
}

func (mock *mockOrderRepo) GetNumOfOrdersPerMonth(ctx context.Context, userId int) (*[]orderEntity.AggregatedOrdersByMonth, error) {
	args := mock.Called(ctx, userId)

	return args.Get(0).(*[]orderEntity.AggregatedOrdersByMonth), args.Error(1)
}

func (mock *mockOrderRepo) GetOrder(ctx context.Context, userId, orderId int) (*orderEntity.Order, error) {
	args := mock.Called(ctx, userId, orderId)

	return args.Get(0).(*orderEntity.Order), args.Error(1)
}

type OrderUsecaseTestSuite struct {
	suite.Suite
}

// func (suite *OrderUsecaseTestSuite) TestCreateOrder() {
// 	tests := []struct {
// 		name string
// 		ctx  context.Context
// 		data *orderEntity.Order
// 	}{}
// }

func TestOrderUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderUsecaseTestSuite))
}
