package webhook

import (
	"context"
	"net/http"
	"oms/database"

	"github.com/gin-gonic/gin"
)

func GetWebhookEvents(c *gin.Context) {
	collection := database.GetGlobalDatabase().GetCollection("webhook_events")
	ctx := context.Background()
	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}
	defer cursor.Close(ctx)

	var events []WebhookEvent
	for cursor.Next(ctx) {
		var event WebhookEvent
		if err := cursor.Decode(&event); err == nil {
			events = append(events, event)
		}
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}
