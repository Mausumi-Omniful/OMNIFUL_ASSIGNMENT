package pubsub

import (
	"context"
)

// Publisher interface defines the behavior of all publishers
type Publisher interface {
	Publish(ctx context.Context, msg *Message) error
	PublishAsync(ctx context.Context, msg *Message) error
	PublishBatch(ctx context.Context, msgs []*Message) error
	Close()
}

// Subscriber interface defines the behavior of all subscribers
type Subscriber interface {
	Subscribe(ctx context.Context) error
	RegisterHandler(topic string, handler IPubSubMessageHandler)
	Close()
}
