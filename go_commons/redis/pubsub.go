package redis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

// Publish publishes a message to the given channel.
// It returns the number of clients that received the message.
func (r *Client) Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	// Assert the underlying client to the redis.UniversalClient interface which has Publish.
	universalClient, ok := r.Cmdable.(redis.UniversalClient)
	if !ok {
		return 0, errors.New("underlying client does not support Publish")
	}
	return universalClient.Publish(ctx, channel, message).Result()
}

// Subscribe subscribes to the given channels and returns a PubSub instance.
// This PubSub instance should be used to receive messages.
func (r *Client) Subscribe(ctx context.Context, channels ...string) (*redis.PubSub, error) {
	// Assert the underlying client to the redis.UniversalClient interface which has Subscribe.
	universalClient, ok := r.Cmdable.(redis.UniversalClient)
	if !ok {
		return nil, errors.New("underlying client does not support Subscribe")
	}
	pubsub := universalClient.Subscribe(ctx, channels...)

	// Optionally, you can call pubsub.ReceiveSubscription(ctx) to ensure that the subscription is ready.
	return pubsub, nil
}

// SubscribeChannel subscribes to the given channels and returns a channel for receiving messages.
// It internally uses the Subscribe method and returns the Go channel from the PubSub instance.
func (r *Client) SubscribeChannel(ctx context.Context, channels ...string) (<-chan *redis.Message, error) {
	// Get a PubSub instance
	pubsub, err := r.Subscribe(ctx, channels...)
	if err != nil {
		return nil, err
	}

	// Optionally, wait for the subscription confirmation.
	// Without this step, you might miss early published messages.
	if _, err := pubsub.Receive(ctx); err != nil {
		return nil, err
	}

	msgChannel := pubsub.Channel()
	return msgChannel, nil
}
