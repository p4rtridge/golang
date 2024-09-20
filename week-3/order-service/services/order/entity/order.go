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
	if order != nil {
		order.Id = id
	}
}

func (order *Order) SetUserId(id int) {
	if order != nil {
		order.UserId = id
	}
}

func (order *Order) SetTotalPrice(price float32) {
	if order != nil {
		order.TotalPrice = price
	}
}

func (order *Order) SetCreatedAt(ca time.Time) {
	if order != nil {
		order.CreatedAt = ca
	}
}

func (order *Order) SetUpdatedAt(ua *time.Time) {
	if order != nil {
		order.UpdatedAt = ua
	}
}

func (order *Order) AddItem(item OrderItem) {
	if order != nil {
		order.Items = append(order.Items, item)
	}
}

func (order *Order) GetIdSafe() int {
	if order != nil {
		return order.Id
	}

	return 0
}

func (order *Order) GetUserIdSafe() int {
	if order != nil {
		return order.UserId
	}

	return 0
}

func (order *Order) GetTotalPriceSafe() float32 {
	if order != nil {
		return order.TotalPrice
	}

	return 0
}

func (order *Order) GetItemsSafe() []OrderItem {
	if order != nil {
		return order.Items
	}

	return []OrderItem{}
}

func (order *Order) GetItemSafe(idx int) *OrderItem {
	if order != nil && idx >= 0 && idx < len(order.Items) {
		return &order.Items[idx]
	}

	return nil
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
	if item != nil {
		item.OrderId = id
	}
}

func (item *OrderItem) SetProductId(productId int) {
	if item != nil {
		item.ProductId = productId
	}
}

func (item *OrderItem) SetProductName(productName string) {
	if item != nil {
		item.ProductName = productName
	}
}

func (item *OrderItem) SetProductPrice(productPrice float32) {
	if item != nil {
		item.ProductPrice = productPrice
	}
}

func (item OrderItem) GetProductId() int {
	return item.ProductId
}

func (item OrderItem) GetQuantity() int {
	return item.Quantity
}

func (item OrderItem) GetProductName() string {
	return item.ProductName
}

func (item OrderItem) GetProductPrice() float32 {
	return item.ProductPrice
}
