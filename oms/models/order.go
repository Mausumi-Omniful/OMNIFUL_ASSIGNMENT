package models

import (
	"fmt"
	"math/rand"
	"time"
)


type OrderStatus string

const (
	OrderStatusOnHold OrderStatus = "on_hold"
	OrderStatusNewOrder OrderStatus = "new_order"
    OrderStatusCancelled OrderStatus = "cancelled"
)




func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusOnHold, OrderStatusNewOrder, OrderStatusCancelled:
		return true
	default:
		return false
	}
}






type Order struct {
	ID        string      `json:"id"`
	SKU       string      `json:"sku"`
	Location  string      `json:"location"`
	TenantID  string      `json:"tenant_id"`
	SellerID  string      `json:"seller_id"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}





// generateOrderID 
func generateOrderID() string {

	timestamp := time.Now().UnixNano()
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomStr := make([]byte, 6)
	for i := range randomStr {
		randomStr[i] = letters[rand.Intn(len(letters))]
	}

	return fmt.Sprintf("ORD-%d-%s", timestamp, string(randomStr))
}




// NewOrder
func NewOrder(sku, location, tenantID, sellerID string) *Order {
	now := time.Now()
	return &Order{
		ID:        generateOrderID(),
		SKU:       sku,
		Location:  location,
		TenantID:  tenantID,
		SellerID:  sellerID,
		Status:    OrderStatusOnHold,
		CreatedAt: now,
		UpdatedAt: now,
	}
}










func (o *Order) IsValid() bool {
	return o.ID != "" && o.SKU != "" && o.Location != "" && o.TenantID != "" && o.SellerID != ""
}

func (o *Order) UpdateStatus(status OrderStatus) {
	o.Status = status
	o.UpdatedAt = time.Now()
}
