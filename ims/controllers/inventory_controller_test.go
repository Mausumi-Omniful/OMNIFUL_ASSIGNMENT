package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Mausumi-Omniful/ims/models"
	"github.com/gin-gonic/gin"
)

func uniqueSuffix() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func TestCreateInventory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/inventories", CreateInventory)

	suffix := uniqueSuffix()
	inv := models.Inventory{ProductID: "P001" + suffix, SKU: "SKU001" + suffix, Location: "A1" + suffix, TenantID: "tenant1", SellerID: "seller1", Quantity: 10}
	jsonBody, _ := json.Marshal(inv)

	req, _ := http.NewRequest("POST", "/inventories", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "Inventory item created successfully" {
		t.Errorf("Expected message 'Inventory item created successfully', got %v", resp["message"])
	}
}








func TestGetInventories(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/inventories", GetInventories)

	req, _ := http.NewRequest("GET", "/inventories", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["data"] == nil {
		t.Errorf("Expected data in response, got %v", resp)
	}
}






func TestUpdateInventory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/inventories", CreateInventory)
	r.PUT("/inventories/:id", UpdateInventory)

	suffix := uniqueSuffix()
	inv := models.Inventory{ProductID: "P002" + suffix, SKU: "SKU002" + suffix, Location: "B1" + suffix, TenantID: "tenant2", SellerID: "seller2", Quantity: 5}
	jsonBody, _ := json.Marshal(inv)
	req, _ := http.NewRequest("POST", "/inventories", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	invID := ""
	if itemData, ok := resp["item"].(map[string]interface{}); ok {
		invID = fmt.Sprintf("%v", itemData["id"])
	}
	if invID == "" {
		t.Fatalf("No inventory ID returned")
	}

	update := map[string]interface{}{"product_id": "P002U" + suffix, "sku": "SKU002U" + suffix, "location": "B2" + suffix, "tenant_id": "tenant2", "seller_id": "seller2", "quantity": 15}
	updateBody, _ := json.Marshal(update)
	updateReq, _ := http.NewRequest("PUT", "/inventories/"+invID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateReq)

	if updateW.Code != 200 {
		t.Errorf("Expected status 200, got %d. Body: %s", updateW.Code, updateW.Body.String())
	}
}









func TestDeleteInventory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/inventories", CreateInventory)
	r.DELETE("/inventories/:id", DeleteInventory)

	suffix := uniqueSuffix()
	inv := models.Inventory{ProductID: "P003" + suffix, SKU: "SKU003" + suffix, Location: "C1" + suffix, TenantID: "tenant3", SellerID: "seller3", Quantity: 7}
	jsonBody, _ := json.Marshal(inv)
	req, _ := http.NewRequest("POST", "/inventories", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	invID := ""
	if itemData, ok := resp["item"].(map[string]interface{}); ok {
		invID = fmt.Sprintf("%v", itemData["id"])
	}
	if invID == "" {
		t.Fatalf("No inventory ID returned")
	}

	deleteReq, _ := http.NewRequest("DELETE", "/inventories/"+invID, nil)
	deleteW := httptest.NewRecorder()
	r.ServeHTTP(deleteW, deleteReq)

	if deleteW.Code != 200 {
		t.Errorf("Expected status 200, got %d. Body: %s", deleteW.Code, deleteW.Body.String())
	}
}
