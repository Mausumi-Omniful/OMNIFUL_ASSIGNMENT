package models

type Hub struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`        // e.g. Bangalore Hub
	Location string `json:"location"`    // e.g. Karnataka
	TenantID string `json:"tenant_id"`   // Each hub belongs to a tenant
	SellerID string `json:"seller_id"`   // Each hub may serve a seller
}
