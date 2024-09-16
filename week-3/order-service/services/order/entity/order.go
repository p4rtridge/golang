package entity

import "time"

type Order struct {
	Id         int        `json:"id"`
	UserId     int        `json:"user_id"`
	TotalPrice float32    `json:"total_price"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	Items      []Item     `json:"items"`
}

func NewOrder(id, userId int, totalPrice float32, items []Item) Order {
	return Order{
		Id:         id,
		UserId:     userId,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
	}
}

func (order *Order) CalculateTotalPrice() {
}

type Item struct {
	OrderId      int     `json:"order_id"`
	ProductId    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float32 `json:"product_price"`
	Quantity     int     `json:"quantity"`
}

func NewItem(orderId, productId int, productName string, productPrice float32, quantity int) Item {
	return Item{
		OrderId:      orderId,
		ProductId:    productId,
		ProductName:  productName,
		ProductPrice: productPrice,
		Quantity:     quantity,
	}
}
