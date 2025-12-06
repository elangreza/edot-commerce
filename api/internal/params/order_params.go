package params

import errs "github.com/elangreza/edot-commerce/api/internal/error"

type AddToCartRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
}

func (a *AddToCartRequest) Validate() error {
	if a.ProductID == "" {
		return errs.ValidationError{Message: "product_id is required"}
	}

	if a.Quantity < 1 {
		return errs.ValidationError{Message: "quantity must be positive"}
	}

	return nil
}

type (
	GetCartItemsResponse struct {
		ProductID string `json:"cart_id"`
		Quantity  int64  `json:"quantity"`
	}

	GetCartResponse struct {
		CartID string                 `json:"cart_id"`
		Items  []GetCartItemsResponse `json:"items"`
	}
)
