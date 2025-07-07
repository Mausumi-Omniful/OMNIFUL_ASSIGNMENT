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

func uniqueSkuSuffix() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func TestCreateSKU(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/skus", CreateSKU)

	suffix := uniqueSkuSuffix()
	sku := models.SKU{Code: "SKU001" + suffix, Name: "Test SKU", Description: "A test SKU", TenantID: "tenant1", SellerID: "seller1"}
	jsonBody, _ := json.Marshal(sku)

	req, _ := http.NewRequest("POST", "/skus", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "SKU created" {
		t.Errorf("Expected message 'SKU created', got %v", resp["message"])
	}
}






func TestGetSKUs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/skus", GetSKUs)

	req, _ := http.NewRequest("GET", "/skus", nil)
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






func TestUpdateSKU(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/skus", CreateSKU)
	r.PUT("/skus/:id", UpdateSKU)

	suffix := uniqueSkuSuffix()
	sku := models.SKU{Code: "SKU002" + suffix, Name: "SKU To Update", Description: "To be updated", TenantID: "tenant2", SellerID: "seller2"}
	jsonBody, _ := json.Marshal(sku)
	req, _ := http.NewRequest("POST", "/skus", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	skuID := ""
	if skuData, ok := resp["sku"].(map[string]interface{}); ok {
		skuID = fmt.Sprintf("%v", skuData["id"])
	}
	if skuID == "" {
		t.Fatalf("No SKU ID returned")
	}

	update := map[string]interface{}{"sku_code": "SKU002U" + suffix, "name": "Updated SKU", "description": "Updated desc", "tenant_id": "tenant2", "seller_id": "seller2"}
	updateBody, _ := json.Marshal(update)
	updateReq, _ := http.NewRequest("PUT", "/skus/"+skuID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateReq)

	if updateW.Code != 200 {
		t.Errorf("Expected status 200, got %d. Body: %s", updateW.Code, updateW.Body.String())
	}
}







func TestDeleteSKU(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/skus", CreateSKU)
	r.DELETE("/skus/:id", DeleteSKU)

	suffix := uniqueSkuSuffix()
	sku := models.SKU{Code: "SKU003" + suffix, Name: "SKU To Delete", Description: "To be deleted", TenantID: "tenant3", SellerID: "seller3"}
	jsonBody, _ := json.Marshal(sku)
	req, _ := http.NewRequest("POST", "/skus", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	skuID := ""
	if skuData, ok := resp["sku"].(map[string]interface{}); ok {
		skuID = fmt.Sprintf("%v", skuData["id"])
	}
	if skuID == "" {
		t.Fatalf("No SKU ID returned")
	}

	deleteReq, _ := http.NewRequest("DELETE", "/skus/"+skuID, nil)
	deleteW := httptest.NewRecorder()
	r.ServeHTTP(deleteW, deleteReq)

	if deleteW.Code != 200 {
		t.Errorf("Expected status 200, got %d. Body: %s", deleteW.Code, deleteW.Body.String())
	}
}
