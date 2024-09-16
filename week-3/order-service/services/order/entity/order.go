package entity

import "time"

type OrderStatus string

const (
	PENDING  OrderStatus = "pending"
	DONE     OrderStatus = "done"
	CANCELED OrderStatus = "canceled"
)

type Order struct {
	Id         int         `json:"id"`
	UserId     int         `json:"user_id"`
	TotalPrice float32     `json:"total_price"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  *time.Time  `json:"updated_at"`
	Items      []OrderItem `json:"items"`
	Status     OrderStatus `json:"order_status"`
}

func NewOrder(id, userId int, totalPrice float32, items []OrderItem) Order {
	return Order{
		Id:         id,
		UserId:     userId,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
		Items:      items,
	}
}

func (order *Order) SetId(id int) {
	order.Id = id
}

func (order *Order) CalculatePrice() float32 {
	totalPrice := float32(0.0)

	for _, item := range order.Items {
		totalPrice += float32(item.Quantity) * item.ProductPrice
	}

	order.TotalPrice = totalPrice

	return totalPrice
}

type OrderItem struct {
	OrderId      int     `json:"order_id"`
	ProductId    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float32 `json:"product_price"`
	Quantity     int     `json:"quantity"`
}

func NewOrderItem(orderId, productId int, productName string, productPrice float32, quantity int) OrderItem {
	return OrderItem{
		OrderId:      orderId,
		ProductId:    productId,
		ProductName:  productName,
		ProductPrice: productPrice,
		Quantity:     quantity,
	}
}

func (item *OrderItem) SetOrderId(id int) {
	item.OrderId = id
}
