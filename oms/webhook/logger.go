package webhook

import (
	"context"
	"encoding/json"
	"oms/database"
	"time"
)

func LogWebhookEvent(ctx context.Context, eventType string, payload interface{}) error {
	data, _ := json.Marshal(payload)
	event := WebhookEvent{
		EventType: eventType,
		Payload:   string(data),
		CreatedAt: time.Now(),
	}
	return database.SaveWebhookEvent(ctx, &event)
}
