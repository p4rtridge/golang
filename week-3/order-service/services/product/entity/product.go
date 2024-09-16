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
	product.Id = id
}

func (product *Product) SetQuantity(q int) {
	product.Quantity = q
}
