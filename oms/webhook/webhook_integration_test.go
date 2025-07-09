package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"oms/database"

	"github.com/gin-gonic/gin"
)

type responseBody struct {
	Events []WebhookEvent `json:"events"`
}

func setupTestDB(t *testing.T) *database.Database {
	mongoURI := "mongodb://myuser:mypassword@localhost:27018/mydb?authSource=admin"
	mongoDBName := "mydb"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := database.NewDatabase(ctx, mongoURI, mongoDBName)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	database.SetGlobalDatabase(db)
	t.Cleanup(func() {
		db.GetCollection("webhook_events").Drop(ctx)
		db.Close(ctx)
		cancel()
	})
	return db
}

func TestWebhookEventsEndpoint_E2E(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Setup DB and insert test data
	db := setupTestDB(t)
	ctx := context.Background()
	testEvent := WebhookEvent{
		EventType: "test_event",
		Payload:   "{\"foo\":\"bar\"}",
		CreatedAt: time.Now(),
	}
	_, err := db.GetCollection("webhook_events").InsertOne(ctx, testEvent)
	if err != nil {
		t.Fatalf("Failed to insert test webhook event: %v", err)
	}

	// Register the real handler
	r.GET("/api/v1/webhook/events", GetWebhookEvents)

	req, _ := http.NewRequest("GET", "/api/v1/webhook/events", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp responseBody
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	found := false
	for _, event := range resp.Events {
		if event.EventType == "test_event" && event.Payload == "{\"foo\":\"bar\"}" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Test event not found in response. Response: %s", w.Body.String())
	}
}
