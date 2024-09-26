package entity

import "time"

type OrderRequest struct {
	Items []ProductItem `json:"items"`
}

type ProductItem struct {
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type AggregatedOrdersByMonth struct {
	Time        time.Time
	NumOfOrders int
}

type OrdersSummarize struct {
	UserId                   int     `json:"user_id"`
	Username                 string  `json:"username"`
	NumOfOrders              int     `json:"num_of_orders"`
	SumOrderPrice            float32 `json:"sum_order_price"`
	AverageOrderItemQuantity float32 `json:"average_order_item_quantity"`
}

func (data OrderRequest) Validate() error {
	if len(data.Items) < 1 {
		return ErrItemEmpty
	}

	for _, item := range data.Items {
		if item.ProductId == 0 {
			return ErrMissingField
		}

		if item.Quantity == 0 {
			return ErrMissingField
		}
	}

	return nil
}

func (data OrderRequest) GetItems() []ProductItem {
	return data.Items
}

func (data ProductItem) GetItemId() int {
	return data.ProductId
}

func (data ProductItem) GetItemQuantity() int {
	return data.Quantity
}
