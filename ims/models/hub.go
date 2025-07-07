package models

import (
	"errors"
	"strings"
)

type Hub struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	TenantID string `json:"tenant_id"`
	SellerID string `json:"seller_id"`
}

// Validate checks if the hub has required fields
func (h *Hub) Validate() error {
	if strings.TrimSpace(h.Name) == "" {
		return errors.New("hub name is required")
	}
	if strings.TrimSpace(h.Location) == "" {
		return errors.New("hub location is required")
	}
	if strings.TrimSpace(h.TenantID) == "" {
		return errors.New("tenant_id is required")
	}
	return nil
}
