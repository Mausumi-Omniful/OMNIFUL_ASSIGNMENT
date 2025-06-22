package models

type Hub struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`      
	Location string `json:"location"`   
	TenantID string `json:"tenant_id"`   
	SellerID string `json:"seller_id"`  
}