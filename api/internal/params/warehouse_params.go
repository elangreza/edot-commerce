package params

type SetWarehouseStatusRequest struct {
	WarehouseID int64 `json:"warehouse_id"`
	IsActive    bool  `json:"is_active"`
}
