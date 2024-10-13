package rpc

import (
	"context"
	"kitchen/services/common/genproto/orders"
	"kitchen/services/orders/entity"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderRepo interface {
	CreateOrder(context.Context, *entity.Order) error
	GetOrders(context.Context, int32) ([]*entity.Order, error)
}

type grpcClient struct {
	client orders.OrderServiceClient
}

func NewGRPCClient(addr string) OrderRepo {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}

	return &grpcClient{
		client: orders.NewOrderServiceClient(conn),
	}
}

func (repo *grpcClient) CreateOrder(ctx context.Context, order *entity.Order) error {
	req := &orders.CreateOrderRequest{
		CustomerID: order.CustomerID,
		ProductID:  order.ProductID,
		Quantity:   order.Quantity,
	}

	_, err := repo.client.CreateOrder(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (repo *grpcClient) GetOrders(ctx context.Context, customerID int32) ([]*entity.Order, error) {
	req := &orders.GetOrdersRequest{
		CustomerID: customerID,
	}

	resp, err := repo.client.GetOrders(ctx, req)
	if err != nil {
		return nil, err
	}

	orders := make([]*entity.Order, 0)

	for _, order := range resp.Orders {
		o := entity.NewOrder(order.OrderID, order.CustomerID, order.ProductID, order.Quantity)

		orders = append(orders, &o)
	}

	return orders, nil
}
