package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/omniful/go_commons/log"
)

// IMSClient handles communication with the IMS service
type IMSClient struct {
	baseURL    string
	httpClient *http.Client
}

// SKU represents a SKU from IMS service
type SKU struct {
	ID          int    `json:"id"`
	Code        string `json:"sku_code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenant_id"`
	SellerID    string `json:"seller_id"`
}

// Hub represents a hub from IMS service
type Hub struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	TenantID string `json:"tenant_id"`
	SellerID string `json:"seller_id"`
}

// Inventory represents inventory from IMS service
type Inventory struct {
	ID        int    `json:"id"`
	ProductID string `json:"product_id"`
	SKU       string `json:"sku"`
	Location  string `json:"location"`
	TenantID  string `json:"tenant_id"`
	SellerID  string `json:"seller_id"`
	Quantity  int    `json:"quantity"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// SKUResponse represents the response from IMS SKU API
type SKUResponse struct {
	Data   []SKU  `json:"data"`
	Source string `json:"source"`
}

// HubResponse represents the response from IMS Hub API
type HubResponse struct {
	Data   []Hub  `json:"data"`
	Source string `json:"source"`
}

// InventoryResponse represents the response from IMS Inventory API
type InventoryResponse struct {
	Data   []Inventory `json:"data"`
	Source string      `json:"source"`
}

// NewIMSClient creates a new IMS API client
func NewIMSClient(baseURL string) *IMSClient {
	return &IMSClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetSKUs fetches all SKUs from the IMS service
func (c *IMSClient) GetSKUs() ([]SKU, error) {
	url := fmt.Sprintf("%s/sku/", c.baseURL)

	log.Infof("Fetching SKUs from IMS: %s", url)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch SKUs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IMS API returned status %d", resp.StatusCode)
	}

	var skuResponse SKUResponse
	if err := json.NewDecoder(resp.Body).Decode(&skuResponse); err != nil {
		return nil, fmt.Errorf("failed to decode SKU response: %w", err)
	}

	log.Infof("Successfully fetched %d SKUs from IMS", len(skuResponse.Data))
	return skuResponse.Data, nil
}

// GetHubs fetches all hubs from the IMS service
func (c *IMSClient) GetHubs() ([]Hub, error) {
	url := fmt.Sprintf("%s/hub/", c.baseURL)

	log.Infof("Fetching Hubs from IMS: %s", url)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch hubs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IMS API returned status %d", resp.StatusCode)
	}

	var hubResponse HubResponse
	if err := json.NewDecoder(resp.Body).Decode(&hubResponse); err != nil {
		return nil, fmt.Errorf("failed to decode hub response: %w", err)
	}

	log.Infof("Successfully fetched %d hubs from IMS", len(hubResponse.Data))
	return hubResponse.Data, nil
}

// GetInventory fetches all inventory from the IMS service
func (c *IMSClient) GetInventory() ([]Inventory, error) {
	url := fmt.Sprintf("%s/inventory/", c.baseURL)

	log.Infof("Fetching Inventory from IMS: %s", url)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch inventory: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IMS API returned status %d", resp.StatusCode)
	}

	var inventoryResponse InventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&inventoryResponse); err != nil {
		return nil, fmt.Errorf("failed to decode inventory response: %w", err)
	}

	log.Infof("Successfully fetched %d inventory items from IMS", len(inventoryResponse.Data))
	return inventoryResponse.Data, nil
}

// ValidateSKU checks if a SKU exists and matches the tenant/seller
func (c *IMSClient) ValidateSKU(skuCode, tenantID, sellerID string) (bool, error) {
	skus, err := c.GetSKUs()
	if err != nil {
		return false, err
	}

	for _, sku := range skus {
		if sku.Code == skuCode {
			// If SKU has empty tenant_id and seller_id, it's available to any tenant/seller
			if (sku.TenantID == "" && sku.SellerID == "") ||
				(sku.TenantID == tenantID && sku.SellerID == sellerID) {
				return true, nil
			}
		}
	}

	return false, nil
}

// ValidateHub checks if a hub exists and matches the tenant/seller
func (c *IMSClient) ValidateHub(hubName, tenantID, sellerID string) (bool, error) {
	hubs, err := c.GetHubs()
	if err != nil {
		return false, err
	}

	for _, hub := range hubs {
		if hub.Name == hubName {
			// If hub has empty tenant_id and seller_id, it's available to any tenant/seller
			if (hub.TenantID == "" && hub.SellerID == "") ||
				(hub.TenantID == tenantID && hub.SellerID == sellerID) {
				return true, nil
			}
		}
	}

	return false, nil
}

// CheckInventoryAvailability checks if inventory is available for a specific SKU at a location
func (c *IMSClient) CheckInventoryAvailability(sku, location, tenantID, sellerID string) (bool, int, error) {
	log.Infof("Checking inventory availability - SKU: %s, Location: %s, Tenant: %s, Seller: %s", sku, location, tenantID, sellerID)

	// Fetch all inventory from IMS
	inventory, err := c.GetInventory()
	if err != nil {
		log.WithError(err).Error("❌ Failed to fetch inventory from IMS")
		return false, 0, fmt.Errorf("failed to fetch inventory: %w", err)
	}

	// Find matching inventory item
	for _, item := range inventory {
		// Check if SKU and location match
		if item.SKU == sku && item.Location == location {
			// Check tenant/seller permissions
			// If inventory has empty tenant_id and seller_id, it's available to any tenant/seller
			if (item.TenantID == "" && item.SellerID == "") ||
				(item.TenantID == tenantID && item.SellerID == sellerID) {

				isAvailable := item.Quantity > 0
				log.Infof("✅ Inventory check result - SKU: %s, Location: %s, Available: %t, Quantity: %d",
					sku, location, isAvailable, item.Quantity)

				return isAvailable, item.Quantity, nil
			}
		}
	}

	// No matching inventory found
	log.Warnf("⚠️ No inventory found - SKU: %s, Location: %s, Tenant: %s, Seller: %s",
		sku, location, tenantID, sellerID)
	return false, 0, nil
}
