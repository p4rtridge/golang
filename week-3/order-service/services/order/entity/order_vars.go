package entity

type OrderRequest struct {
	Items []ProductItem `json:"items"`
}

type ProductItem struct {
	ProductId *int `json:"product_id"`
	Quantity  *int `json:"quantity"`
}

func (data OrderRequest) Validate() error {
	if len(data.Items) < 1 {
		return ErrItemEmpty
	}

	for _, item := range data.Items {
		if item.ProductId == nil {
			return ErrMissingField
		}

		if item.Quantity == nil {
			return ErrMissingField
		}
	}

	return nil
}

func (data OrderRequest) GetItems() []ProductItem {
	return data.Items
}

func (data ProductItem) GetItemId() int {
	if data.ProductId != nil {
		return *data.ProductId
	}
	return 0
}

func (data ProductItem) GetItemQuantity() int {
	if data.Quantity != nil {
		return *data.Quantity
	}
	return 0
}
