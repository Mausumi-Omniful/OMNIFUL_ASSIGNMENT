package webhook

import "time"

type WebhookEvent struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	EventType string    `json:"event_type" bson:"event_type"`
	Payload   string    `json:"payload" bson:"payload"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
