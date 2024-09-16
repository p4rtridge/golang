package entity

type ProductItem struct {
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type OrderRequest struct {
	Items []ProductItem `json:"items"`
}

func (data *OrderRequest) Validate() error {
	if len(data.Items) < 1 {
		return ErrItemEmpty
	}

	return nil
}
