package models

import (
	"fmt"
	"math/rand"
	"time"
)

// OrderStatus represents the current status of an order
type OrderStatus string

const (
	// OrderStatusOnHold - Initial state when order is created from CSV
	OrderStatusOnHold OrderStatus = "on_hold"

	// OrderStatusNewOrder - Order is confirmed and inventory is allocated
	OrderStatusNewOrder OrderStatus = "new_order"

	// OrderStatusCancelled - Order is cancelled due to validation failure
	OrderStatusCancelled OrderStatus = "cancelled"
)

// IsValid checks if the order status is valid
func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusOnHold, OrderStatusNewOrder, OrderStatusCancelled:
		return true
	default:
		return false
	}
}

// Order represents an order in the system
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

// generateOrderID creates a unique order ID
func generateOrderID() string {
	// Use timestamp for uniqueness
	timestamp := time.Now().UnixNano()

	// Add random string for additional uniqueness
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomStr := make([]byte, 6)
	for i := range randomStr {
		randomStr[i] = letters[rand.Intn(len(letters))]
	}

	return fmt.Sprintf("ORD-%d-%s", timestamp, string(randomStr))
}

// NewOrder creates a new order with default values
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

// IsValid checks if the order has all required fields
func (o *Order) IsValid() bool {
	return o.ID != "" && o.SKU != "" && o.Location != "" && o.TenantID != "" && o.SellerID != ""
}

// UpdateStatus updates the order status and timestamp
func (o *Order) UpdateStatus(status OrderStatus) {
	o.Status = status
	o.UpdatedAt = time.Now()
}
