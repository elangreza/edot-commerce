package params

type AddToCartRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
}
