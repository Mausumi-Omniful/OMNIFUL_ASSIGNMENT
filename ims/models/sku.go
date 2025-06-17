package models

type SKU struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"unique;not null" json:"code"` // e.g. SKU001
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenant_id"`  // Helps to isolate per client
	SellerID    string `json:"seller_id"`  // Each seller has their SKUs
}
