package dchannel

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/omniful/go_commons/log"
	oredis "github.com/omniful/go_commons/redis"
)

const channelCloseMessage = "CHAN_CLOSE"

func NewRedisChannel(name string, redis *oredis.Client) DChannel {
	return &RedisChannel{
		name:  name,
		redis: redis,
	}
}

type RedisChannel struct {
	name     string
	redis    *oredis.Client
	isClosed bool
}

func (c *RedisChannel) Push(ctx context.Context, m string) error {
	if c.isClosed {
		return ErrChanClosed
	}
	_, err := c.redis.Publish(ctx, c.name, m)
	return err
}

func (c *RedisChannel) Listen(ctx context.Context) (string, error) {
	if c.isClosed {
		return "", ErrChanClosed
	}

	pubSub, err := c.redis.Subscribe(ctx, c.name)
	if err != nil {
		return "", err
	}
	defer c.closePubSub(pubSub)
	defer c.unsubscribe(ctx, pubSub)

	// wait for the subscription confirmation. sent by redis when successfully subscribed
	if _, err = pubSub.Receive(ctx); err != nil {
		return "", err
	}

	select {
	case m := <-pubSub.Channel():
		if m.Payload == channelCloseMessage {
			return "", ErrChanClosed
		}
		return m.Payload, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (c *RedisChannel) IsClosed() bool {
	return c.isClosed
}

func (c *RedisChannel) Close() {
	c.Push(context.Background(), channelCloseMessage)
	c.isClosed = true
}

func (c *RedisChannel) unsubscribe(ctx context.Context, pubSub *redis.PubSub) {
	err := pubSub.Unsubscribe(ctx, c.name)
	if err != nil {
		log.Errorf("error while unsubscribing channel")
	}
}

func (c *RedisChannel) closePubSub(pubSub *redis.PubSub) {
	err := pubSub.Close()
	if err != nil {
		log.Errorf("error while closing pubsub client")
	}
}
