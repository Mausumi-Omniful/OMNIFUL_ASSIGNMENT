package models

type SKU struct {
	ID          uint   `json:"id"`
	Code        string `json:"sku_code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenant_id"`
	SellerID    string `json:"seller_id"`
}
