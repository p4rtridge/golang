package entity

import (
	"net/http"
)

type ProductRequest struct {
	Name     string  `json:"name"`
	Image    []byte  `json:"image"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"price"`
}

func (product *ProductRequest) Validate() error {
	mimeType := http.DetectContentType(product.Image)

	switch mimeType {
	case "image/jpeg", "image/png":
		break
	default:
		return ErrInvalidRequestBody
	}

	return nil
}
