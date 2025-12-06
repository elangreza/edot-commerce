package entity

import "github.com/google/uuid"

type Stock struct {
	ID        int64     `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
}

type ReserveStock struct {
	Stocks  []Stock   `json:"stocks"`
	OrderID string    `json:"order_id"`
	UserID  uuid.UUID `json:"user_id"`
}

type ReleaseStock struct {
	// ReservedStockIDs []int64   `json:"reserved_stock_id"`
	OrderID string    `json:"order_id"`
	UserID  uuid.UUID `json:"user_id"`
}

type ConfirmStock struct {
	ReservedStockIDs []int64   `json:"reserved_stock_id"`
	UserID           uuid.UUID `json:"user_id"`
}
