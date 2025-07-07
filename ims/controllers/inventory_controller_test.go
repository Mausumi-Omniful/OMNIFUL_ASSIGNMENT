package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type Inventory struct {
	ID        int    `json:"id,omitempty"`
	ProductID string `json:"product_id"`
	SKU       string `json:"sku"`
	Location  string `json:"location"`
	TenantID  string `json:"tenant_id"`
	SellerID  string `json:"seller_id"`
	Quantity  int    `json:"quantity"`
}

type mockInventoryStore struct{}

func (m *mockInventoryStore) CreateInventory(inv *Inventory) error { return nil }
func (m *mockInventoryStore) GetInventories() ([]Inventory, error) {
	return []Inventory{{ID: 1, ProductID: "P001", SKU: "SKU001", Location: "A1", TenantID: "tenant1", SellerID: "seller1", Quantity: 10}}, nil
}
func (m *mockInventoryStore) UpdateInventory(id string, inv *Inventory) error { return nil }
func (m *mockInventoryStore) DeleteInventory(id string) error                 { return nil }

func TestInventoryHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &mockInventoryStore{}
	r := gin.Default()
	r.POST("/inventories", func(c *gin.Context) {
		var inv Inventory
		if err := c.ShouldBindJSON(&inv); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		if err := store.CreateInventory(&inv); err != nil {
			c.JSON(500, gin.H{"error": "DB error"})
			return
		}
		c.JSON(200, gin.H{"message": "Inventory item created", "item": inv})
	})
	r.GET("/inventories", func(c *gin.Context) {
		invs, _ := store.GetInventories()
		c.JSON(200, gin.H{"data": invs})
	})
	r.PUT("/inventories/:id", func(c *gin.Context) {
		id := c.Param("id")
		var inv Inventory
		if err := c.ShouldBindJSON(&inv); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		if err := store.UpdateInventory(id, &inv); err != nil {
			c.JSON(500, gin.H{"error": "DB error"})
			return
		}
		c.JSON(200, gin.H{"message": "Inventory updated"})
	})
	r.DELETE("/inventories/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := store.DeleteInventory(id); err != nil {
			c.JSON(500, gin.H{"error": "DB error"})
			return
		}
		c.JSON(200, gin.H{"message": "Inventory deleted"})
	})

	tests := []struct {
		name       string
		method     string
		url        string
		body       interface{}
		wantStatus int
		wantField  string
		wantValue  interface{}
	}{
		{"create inventory", "POST", "/inventories", Inventory{ProductID: "P001", SKU: "SKU001", Location: "A1", TenantID: "tenant1", SellerID: "seller1", Quantity: 10}, 200, "message", "Inventory item created"},
		{"get inventories", "GET", "/inventories", nil, 200, "data", nil},
		{"update inventory", "PUT", "/inventories/1", Inventory{ProductID: "P002", SKU: "SKU002", Location: "B1", TenantID: "tenant2", SellerID: "seller2", Quantity: 5}, 200, "message", "Inventory updated"},
		{"delete inventory", "DELETE", "/inventories/1", nil, 200, "message", "Inventory deleted"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				jsonBody, _ := json.Marshal(tt.body)
				req, _ = http.NewRequest(tt.method, tt.url, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, _ = http.NewRequest(tt.method, tt.url, nil)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("%s: expected status %d, got %d. Body: %s", tt.name, tt.wantStatus, w.Code, w.Body.String())
			}
			if tt.wantField != "" {
				var resp map[string]interface{}
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				if tt.wantValue != nil && resp[tt.wantField] != tt.wantValue {
					t.Errorf("%s: expected %s = %v, got %v", tt.name, tt.wantField, tt.wantValue, resp[tt.wantField])
				}
				if tt.wantField == "data" && resp[tt.wantField] == nil {
					t.Errorf("%s: expected data in response, got %v", tt.name, resp)
				}
			}
		})
	}
}
