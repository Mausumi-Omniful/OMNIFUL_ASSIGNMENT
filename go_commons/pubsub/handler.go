package pubsub

import "context"

type IPubSubMessageHandler interface {
	Process(ctx context.Context, message *Message) error
}
