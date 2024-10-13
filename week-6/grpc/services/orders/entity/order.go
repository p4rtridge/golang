package entity

import "kitchen/services/common/genproto/orders"

type Order struct {
	OrderID    int32 `json:"order_id,omitempty"`
	CustomerID int32 `json:"customer_id,omitempty"`
	ProductID  int32 `json:"product_id,omitempty"`
	Quantity   int32 `json:"quantity,omitempty"`
}

func NewOrder(orderID int32, customerID int32, productID int32, quantity int32) Order {
	return Order{
		OrderID:    orderID,
		CustomerID: customerID,
		ProductID:  productID,
		Quantity:   quantity,
	}
}

type OrderFactory struct{}

func (OrderFactory) CreateFromProto(protoOrder *orders.CreateOrderRequest) Order {
	return Order{
		OrderID:    1,
		CustomerID: protoOrder.CustomerID,
		ProductID:  protoOrder.ProductID,
		Quantity:   protoOrder.Quantity,
	}
}
