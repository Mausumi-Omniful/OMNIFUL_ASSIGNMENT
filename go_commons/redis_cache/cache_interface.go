package cache

import (
	"context"
	"time"
)

type ICache interface {
	Get(ctx context.Context, key string, data interface{}) (bool, error)
	Set(ctx context.Context, key string, data interface{}, ttl time.Duration) (bool, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, k string, expiry time.Duration) (bool, error)
	PipedMSet(ctx context.Context, kvArr []KVIn, d time.Duration) error
	PipedMGet(ctx context.Context, kvArr []*KVOut) (okCount int, err error)
	MSet(ctx context.Context, kvArr []KVIn) error
	ZCard(ctx context.Context, key string) (int64, error)
	ZRange(ctx context.Context, key string, start int64, end int64) ([]string, error)
	ZRangeByScore(ctx context.Context, key string, min string, max string) ([]string, error)
	ZRem(ctx context.Context, key string, members []string) (int64, error)
	ZAdd(ctx context.Context, key string, items []SortedSetItem) (int64, error)
	Unlink(ctx context.Context, keys []string) (int64, error)
	HSet(ctx context.Context, key string, field, value string) (int64, error)
	HSetAll(ctx context.Context, key string, values map[string]interface{}) (int64, error)
	HGet(ctx context.Context, key string, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields []string) (int64, error)
	HIncrBy(ctx context.Context, key string, field string, incr int64) (int64, error)
	HIncrByFloat(ctx context.Context, key string, field string, incr float64) (float64, error)
	IncrBy(ctx context.Context, key string, value int64) error
	DecrBy(ctx context.Context, key string, value int64) error
	MGet(ctx context.Context, keys []string) ([]interface{}, error)
	IncrWithLimit(ctx context.Context, key string, limit int, expiry time.Duration) (int64, error)
	IncrWithMOD(ctx context.Context, key string, mod uint64, expiry time.Duration) (uint64, error)
}
