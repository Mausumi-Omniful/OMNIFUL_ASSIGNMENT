package models

type SKU struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"unique;not null" json:"code"`
	Name        string `json:"name"`
	SKUCode     string `json:"sku_code"`
	Description string `json:"description"`
	TenantID    string `json:"tenant_id"`
	SellerID    string `json:"seller_id"`
}
