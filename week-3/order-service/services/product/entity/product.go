package entity

import "time"

type Product struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Quantity  int        `json:"quantity"`
	Price     float32    `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewProduct(id int, name string, quantity int, price float32) Product {
	return Product{
		Id:        id,
		Name:      name,
		Quantity:  quantity,
		Price:     price,
		CreatedAt: time.Now(),
		UpdatedAt: &time.Time{},
	}
}

func (product *Product) SetId(id int) {
	if product != nil {
		product.Id = id
	}
}

func (product *Product) SetQuantity(q int) {
	if product != nil {
		product.Quantity = q
	}
}

func (product Product) GetId() int {
	return product.Id
}

func (product Product) GetName() string {
	return product.Name
}

func (product Product) GetQuantity() int {
	return product.Quantity
}

func (product Product) GetPrice() float32 {
	return product.Price
}
