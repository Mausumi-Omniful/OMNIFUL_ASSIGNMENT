package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type IMSClient struct {
	baseURL    string
	httpClient *http.Client
}

type SKU struct {
	ID          int    `json:"id"`
	Code        string `json:"sku_code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenant_id"`
	SellerID    string `json:"seller_id"`
}

type Hub struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	TenantID string `json:"tenant_id"`
	SellerID string `json:"seller_id"`
}

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

type SKUResponse struct {
	Data   []SKU  `json:"data"`
	Source string `json:"source"`
}

type HubResponse struct {
	Data   []Hub  `json:"data"`
	Source string `json:"source"`
}

type InventoryResponse struct {
	Data   []Inventory `json:"data"`
	Source string      `json:"source"`
}

func NewIMSClient(baseURL string) *IMSClient {
	return &IMSClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}



// getskus
func (c *IMSClient) GetSKUs() ([]SKU, error) {
	url := fmt.Sprintf("%s/sku/", c.baseURL)
	fmt.Printf("Fetching SKUs from IMS: %s\n", url)

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

	fmt.Printf("Successfully fetched %d SKUs from IMS\n", len(skuResponse.Data))
	return skuResponse.Data, nil
}



// gethubs
func (c *IMSClient) GetHubs() ([]Hub, error) {
	url := fmt.Sprintf("%s/hub/", c.baseURL)
	fmt.Printf("Fetching Hubs from IMS: %s\n", url)

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

	fmt.Printf("Successfully fetched %d hubs from IMS\n", len(hubResponse.Data))
	return hubResponse.Data, nil
}



// getinventory
func (c *IMSClient) GetInventory() ([]Inventory, error) {
	url := fmt.Sprintf("%s/inventory/", c.baseURL)
	fmt.Printf("Fetching Inventory from IMS: %s\n", url)

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

	fmt.Printf("Successfully fetched %d inventory items from IMS\n", len(inventoryResponse.Data))
	return inventoryResponse.Data, nil
}





// validatesku
func (c *IMSClient) ValidateSKU(skuCode, tenantID, sellerID string) (bool, error) {
	skus, err := c.GetSKUs()
	if err != nil {
		return false, err
	}

	for _, sku := range skus {
		if sku.Code == skuCode {
			if (sku.TenantID == "" && sku.SellerID == "") ||
				(sku.TenantID == tenantID && sku.SellerID == sellerID) {
				return true, nil
			}
		}
	}

	return false, nil
}




// validatehub
func (c *IMSClient) ValidateHub(hubName, tenantID, sellerID string) (bool, error) {
	hubs, err := c.GetHubs()
	if err != nil {
		return false, err
	}

	for _, hub := range hubs {
		if hub.Name == hubName {
			if (hub.TenantID == "" && hub.SellerID == "") ||
				(hub.TenantID == tenantID && hub.SellerID == sellerID) {
				return true, nil
			}
		}
	}

	return false, nil
}










// checkinventory
func (c *IMSClient) CheckInventoryAvailability(sku, location, tenantID, sellerID string) (bool, int, error) {
	fmt.Printf("Checking inventory availability - SKU: %s, Location: %s, Tenant: %s, Seller: %s\n",
		sku, location, tenantID, sellerID)

	inventory, err := c.GetInventory()
	if err != nil {
		return false, 0, fmt.Errorf("failed to fetch inventory: %w", err)
	}

	for _, item := range inventory {
		if item.SKU == sku && item.Location == location {
			if (item.TenantID == "" && item.SellerID == "") ||
				(item.TenantID == tenantID && item.SellerID == sellerID) {

				isAvailable := item.Quantity > 0
				fmt.Printf("Inventory check result - SKU: %s, Location: %s, Available: %t, Quantity: %d\n",
					sku, location, isAvailable, item.Quantity)

				return isAvailable, item.Quantity, nil
			}
		}
	}

	fmt.Printf("No inventory found - SKU: %s, Location: %s, Tenant: %s, Seller: %s\n",
		sku, location, tenantID, sellerID)

	return false, 0, nil
}
