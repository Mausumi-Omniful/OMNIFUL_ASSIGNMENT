// Distributed Channel
package dchannel

import (
	"context"
	"github.com/omniful/go_commons/redis"
)

type DChannel interface {
	// Push message to the channel
	Push(ctx context.Context, m string) error
	// Listen for the message from the channel
	Listen(ctx context.Context) (string, error)
	// IsClosed returns true if channel is closed
	IsClosed() bool
	// Close the channel for all listeners
	Close()
}

func New(name string, redis *redis.Client) DChannel {
	return NewRedisChannel(name, redis)
}
