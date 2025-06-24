package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Cmder = redis.Cmder

type Pipeliner interface {
	redis.Pipeliner
	ExecContext(context.Context) ([]Cmder, error)
}

type Pipeline struct {
	redis.Pipeliner
}

type StringCmd struct {
	*redis.StringCmd
}

type BoolCmd struct {
	*redis.BoolCmd
}

type IntCmd struct {
	*redis.IntCmd
}

type SliceCmd struct {
	*redis.SliceCmd
}

type StringStringMapCmd struct {
	*redis.StringStringMapCmd
}

type StatusCmd struct {
	*redis.StatusCmd
}

func (c *Pipeline) ExecContext(ctx context.Context) ([]Cmder, error) {
	return c.Pipeliner.Exec(ctx)
}

func (r *Client) Pipeline() Pipeliner {
	return &Pipeline{Pipeliner: r.Cmdable.Pipeline()}
}
