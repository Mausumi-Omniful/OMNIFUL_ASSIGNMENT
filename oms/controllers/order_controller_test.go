package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/gin-gonic/gin"

	"encoding/json"
	"time"

	"oms/database"
	"oms/models"
	"oms/utils"
)

func TestOrderController_RouteSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/orders/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	req, _ := http.NewRequest("GET", "/orders/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGetOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	
	mongoURI := "mongodb://myuser:mypassword@localhost:27018/mydb?authSource=admin"
	mongoDBName := "mydb"
	s3Bucket := "order-csv-bucket"
	s3Endpoint := "http://localhost:4566"
	awsRegion := "us-east-1"
	sqsEndpoint := "http://localhost:4566"
	sqsQueue := "CreateBulkOrder"

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := database.NewDatabase(ctx, mongoURI, mongoDBName)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Close(ctx)

	// Create OrderRepository
	orderRepo := database.NewOrderRepository(db)

	// Insert a test order
	testOrder := models.NewOrder("sku-test", "loc-test", "tenant-test", "seller-test")
	if err := orderRepo.SaveOrder(ctx, testOrder); err != nil {
		t.Fatalf("Failed to insert test order: %v", err)
	}

	// S3 and SQS 
	s3Uploader, err := utils.NewS3Uploader(s3Bucket, s3Endpoint, awsRegion)
	if err != nil {
		t.Fatalf("Failed to create S3Uploader: %v", err)
	}
	sqsPublisher, err := utils.NewSQSPublisher(sqsQueue, sqsEndpoint, awsRegion)
	if err != nil {
		t.Fatalf("Failed to create SQSPublisher: %v", err)
	}

	// Create the real controller
	orderController := &OrderController{
		S3Uploader:   s3Uploader,
		SQSPublisher: sqsPublisher,
		OrderRepo:    orderRepo,
	}

	// Register the real endpoint
	r.GET("/orders", orderController.ListOrders)

	// Call the endpoint
	req, _ := http.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Check that the test order is in the response
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	orders, ok := resp["orders"].([]interface{})
	if !ok {
		
		t.Fatalf("Response does not contain 'orders' field: %v", resp)
	}
	found := false
	for _, o := range orders {
		orderMap, ok := o.(map[string]interface{})
		if ok &&
			orderMap["sku"] == testOrder.SKU &&
			orderMap["location"] == testOrder.Location &&
			orderMap["status"] == string(testOrder.Status) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Test order not found in response. Response: %s", w.Body.String())
	}
}
