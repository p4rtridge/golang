package entity

import "time"

type OrderStatus string

type Order struct {
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  *time.Time  `json:"updated_at"`
	Items      []OrderItem `json:"items"`
	Id         int         `json:"id"`
	UserId     int         `json:"user_id"`
	TotalPrice float32     `json:"total_price"`
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

func (order *Order) SetTotalPrice(price float32) {
	order.TotalPrice = price
}

type OrderItem struct {
	ProductName  string  `json:"product_name"`
	OrderId      int     `json:"order_id"`
	ProductId    int     `json:"product_id"`
	Quantity     int     `json:"quantity"`
	ProductPrice float32 `json:"product_price"`
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

func (item *OrderItem) SetProductId(productId int) {
	item.ProductId = productId
}

func (item *OrderItem) SetProductName(productName string) {
	item.ProductName = productName
}

func (item *OrderItem) SetProductPrice(productPrice float32) {
	item.ProductPrice = productPrice
}
