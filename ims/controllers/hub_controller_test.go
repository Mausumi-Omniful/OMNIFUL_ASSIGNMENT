package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/models"
	"github.com/Mausumi-Omniful/ims/redisclient"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/config"
)

func TestMain(m *testing.M) {
	dir, _ := os.Getwd()
	fmt.Println("Current working directory:", dir)
	os.Chdir("D:/Omniful-Assignment/ims")

	if os.Getenv("CONFIG_SOURCE") == "" {
		os.Setenv("CONFIG_SOURCE", "local")
	}

	err := config.Init(15 * time.Second)
	if err != nil {
		panic("Failed to initialize config: " + err.Error())
	}

	ctx, err := config.TODOContext()
	if err != nil {
		panic("Failed to get config context: " + err.Error())
	}

	err = db.InitPostgres(ctx)
	if err != nil {
		panic("Postgres connection failed: " + err.Error())
	}

	err = redisclient.InitRedis(ctx)
	if err != nil {
		panic("Redis initialization failed: " + err.Error())
	}

	os.Exit(m.Run())
}





func TestCreateHub(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/hubs", CreateHub)

	hub := models.Hub{
		Name:     "Test Hub",
		Location: "NY",
		TenantID: "tenant1",
		SellerID: "seller1",
	}
	jsonBody, _ := json.Marshal(hub)

	req, _ := http.NewRequest("POST", "/hubs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "Hub created" {
		t.Errorf("Expected message 'Hub created', got %v", resp["message"])
	}
}





func TestGetHubs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/hubs", GetHubs)

	req, _ := http.NewRequest("GET", "/hubs", nil)
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





func TestUpdateHub(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/hubs", CreateHub)
	r.PUT("/hubs/:id", UpdateHub)

	hub := models.Hub{Name: "HubToUpdate", Location: "LA", TenantID: "tenant2", SellerID: "seller2"}
	jsonBody, _ := json.Marshal(hub)
	req, _ := http.NewRequest("POST", "/hubs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	hubID := ""
	if hubData, ok := resp["hub"].(map[string]interface{}); ok {
		hubID = fmt.Sprintf("%v", hubData["id"])
	}
	if hubID == "" {
		t.Fatalf("No hub ID returned")
	}

	update := map[string]interface{}{"name": "Updated Hub", "location": "SF", "tenant_id": "tenant2", "seller_id": "seller2"}
	updateBody, _ := json.Marshal(update)
	updateReq, _ := http.NewRequest("PUT", "/hubs/"+hubID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateReq)

	if updateW.Code != 200 {
		t.Errorf("Expected status 200, got %d. Body: %s", updateW.Code, updateW.Body.String())
	}
}










func TestDeleteHub(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/hubs", CreateHub)
	r.DELETE("/hubs/:id", DeleteHub)

	hub := models.Hub{Name: "HubToDelete", Location: "TX", TenantID: "tenant3", SellerID: "seller3"}
	jsonBody, _ := json.Marshal(hub)
	req, _ := http.NewRequest("POST", "/hubs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	hubID := ""
	if hubData, ok := resp["hub"].(map[string]interface{}); ok {
		hubID = fmt.Sprintf("%v", hubData["id"])
	}
	if hubID == "" {
		t.Fatalf("No hub ID returned")
	}

	deleteReq, _ := http.NewRequest("DELETE", "/hubs/"+hubID, nil)
	deleteW := httptest.NewRecorder()
	r.ServeHTTP(deleteW, deleteReq)

	if deleteW.Code != 200 {
		t.Errorf("Expected status 200, got %d. Body: %s", deleteW.Code, deleteW.Body.String())
	}
}
