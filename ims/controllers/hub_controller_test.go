package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mausumi-ghadei-omniful/ims/models"
)

type mockHubStore struct{}

func (m *mockHubStore) CreateHub(hub *models.Hub) error { return nil }
func (m *mockHubStore) GetHubs() ([]models.Hub, error) {
	return []models.Hub{{ID: 1, Name: "Test Hub", Location: "NY", TenantID: "tenant1", SellerID: "seller1"}}, nil
}
func (m *mockHubStore) UpdateHub(id string, hub *models.Hub) error { return nil }
func (m *mockHubStore) DeleteHub(id string) error                  { return nil }




func TestHubHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &mockHubStore{}
	r := gin.Default()
	r.POST("/hubs", func(c *gin.Context) {
		var hub models.Hub
		if err := c.ShouldBindJSON(&hub); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		if err := store.CreateHub(&hub); err != nil {
			c.JSON(500, gin.H{"error": "DB error"})
			return
		}
		c.JSON(200, gin.H{"message": "Hub created", "hub": hub})
	})
	r.GET("/hubs", func(c *gin.Context) {
		hubs, _ := store.GetHubs()
		c.JSON(200, gin.H{"data": hubs})
	})
	r.PUT("/hubs/:id", func(c *gin.Context) {
		id := c.Param("id")
		var hub models.Hub
		if err := c.ShouldBindJSON(&hub); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		if err := store.UpdateHub(id, &hub); err != nil {
			c.JSON(500, gin.H{"error": "DB error"})
			return
		}
		c.JSON(200, gin.H{"message": "Hub updated"})
	})
	r.DELETE("/hubs/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := store.DeleteHub(id); err != nil {
			c.JSON(500, gin.H{"error": "DB error"})
			return
		}
		c.JSON(200, gin.H{"message": "Hub deleted"})
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
		{"create hub", "POST", "/hubs", models.Hub{Name: "Test Hub", Location: "NY", TenantID: "tenant1", SellerID: "seller1"}, 200, "message", "Hub created"},
		{"get hubs", "GET", "/hubs", nil, 200, "data", nil},
		{"update hub", "PUT", "/hubs/1", models.Hub{Name: "Updated Hub", Location: "SF", TenantID: "tenant2", SellerID: "seller2"}, 200, "message", "Hub updated"},
		{"delete hub", "DELETE", "/hubs/1", nil, 200, "message", "Hub deleted"},
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
