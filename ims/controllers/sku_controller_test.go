package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type SKU struct {
	ID          int    `json:"id,omitempty"`
	Code        string `json:"sku_code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenant_id"`
	SellerID    string `json:"seller_id"`
}

type mockSKUStore struct{}

func (m *mockSKUStore) CreateSKU(sku *SKU) error { return nil }
func (m *mockSKUStore) GetSKUs() ([]SKU, error) {
	return []SKU{{ID: 1, Code: "SKU001", Name: "Test SKU", Description: "desc", TenantID: "tenant1", SellerID: "seller1"}}, nil
}
func (m *mockSKUStore) UpdateSKU(id string, sku *SKU) error { return nil }
func (m *mockSKUStore) DeleteSKU(id string) error           { return nil }

func TestSKUHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &mockSKUStore{}
	r := gin.Default()
	r.POST("/skus", func(c *gin.Context) {
		var sku SKU
		if err := c.ShouldBindJSON(&sku); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		if err := store.CreateSKU(&sku); err != nil {
			c.JSON(500, gin.H{"error": "DB error"})
			return
		}
		c.JSON(200, gin.H{"message": "SKU created", "sku": sku})
	})
	r.GET("/skus", func(c *gin.Context) {
		skUs, _ := store.GetSKUs()
		c.JSON(200, gin.H{"data": skUs})
	})
	r.PUT("/skus/:id", func(c *gin.Context) {
		id := c.Param("id")
		var sku SKU
		if err := c.ShouldBindJSON(&sku); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		if err := store.UpdateSKU(id, &sku); err != nil {
			c.JSON(500, gin.H{"error": "DB error"})
			return
		}
		c.JSON(200, gin.H{"message": "SKU updated"})
	})
	r.DELETE("/skus/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := store.DeleteSKU(id); err != nil {
			c.JSON(500, gin.H{"error": "DB error"})
			return
		}
		c.JSON(200, gin.H{"message": "SKU deleted"})
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
		{"create sku", "POST", "/skus", SKU{Code: "SKU001", Name: "Test SKU", Description: "desc", TenantID: "tenant1", SellerID: "seller1"}, 200, "message", "SKU created"},
		{"get skus", "GET", "/skus", nil, 200, "data", nil},
		{"update sku", "PUT", "/skus/1", SKU{Code: "SKU002", Name: "Updated SKU", Description: "desc2", TenantID: "tenant2", SellerID: "seller2"}, 200, "message", "SKU updated"},
		{"delete sku", "DELETE", "/skus/1", nil, 200, "message", "SKU deleted"},
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
