package entity

type OrderRequest struct {
	Items []Item `json:"items"`
}

func (data *OrderRequest) Validate() error {
	if len(data.Items) < 1 {
		return ErrItemEmpty
	}

	return nil
}
