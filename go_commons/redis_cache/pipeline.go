package cache

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/redis"
	"time"
)

type Pipeline struct {
	pipeliner  redis.Pipeliner
	serializer ISerializer
}

func (p *Pipeline) Close() error {
	return p.pipeliner.Close()
}

func (p *Pipeline) Set(ctx context.Context, key string, val interface{}, d time.Duration) (*StatusCmd, error) {
	bytes, err := p.serializer.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("serialize - %s", err.Error())
	}

	cmd := p.pipeliner.Set(ctx, key, string(bytes), d)
	return &StatusCmd{statusCmd: &redis.StatusCmd{StatusCmd: cmd}}, nil
}

func (p *Pipeline) Get(ctx context.Context, key string) (*StringCmd, error) {
	cmd := p.pipeliner.Get(ctx, key)
	return &StringCmd{stringCmd: &redis.StringCmd{StringCmd: cmd}, serializer: p.serializer}, nil
}

func (p *Pipeline) Exec(ctx context.Context) error {
	_, err := p.pipeliner.ExecContext(ctx)
	return err
}

type StatusCmd struct {
	statusCmd *redis.StatusCmd
}

func (s *StatusCmd) Result() (string, error) {
	return s.statusCmd.Result()
}

type StringCmd struct {
	stringCmd  *redis.StringCmd
	serializer ISerializer
}

func (s *StringCmd) Result(val interface{}) (ok bool, err error) {
	bytes, err := s.stringCmd.Bytes()
	if err != nil {
		return false, err
	}

	if err = s.serializer.Unmarshal(bytes, val); err != nil {
		return false, fmt.Errorf("deserialize - %s", err.Error())
	}

	return true, nil
}
