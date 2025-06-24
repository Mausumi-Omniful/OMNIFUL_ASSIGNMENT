package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// ZMember represents an item in a sorted set
type ZMember struct {
	Member interface{}
	Score  float64
	Key    string
}

// BZPopMin removes and returns the member with the lowest score from one or more sorted sets,
// or block until one is available
func (r *Client) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) (ZMember, error) {
	ZK, err := r.Cmdable.BZPopMin(ctx, timeout, keys...).Result()
	Z := ZMember{
		Member: ZK.Member,
		Score:  ZK.Score,
		Key:    ZK.Key,
	}

	return Z, err
}

// BZPopMax removes and returns the member with the highest score from one or more sorted sets,
// or block until one is available
func (r *Client) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) (ZMember, error) {
	ZK, err := r.Cmdable.BZPopMax(ctx, timeout, keys...).Result()
	if err != nil {
		return ZMember{}, err
	}

	Z := ZMember{
		Member: ZK.Member,
		Score:  ZK.Score,
		Key:    ZK.Key,
	}

	return Z, err
}

// ZAdd adds one or more members to a sorted set, or update its score if it already exists
func (r *Client) ZAdd(ctx context.Context, key string, members ...ZMember) (int64, error) {
	Z := make([]*redis.Z, 0)
	for _, v := range members {
		Z = append(Z, &redis.Z{
			Score:  v.Score,
			Member: v.Member,
		})
	}

	return r.Cmdable.ZAdd(ctx, key, Z...).Result()
}

// ZAddNX adds only new members, does not update score for already existing members
func (r *Client) ZAddNX(ctx context.Context, key string, members ...ZMember) (int64, error) {
	Z := make([]*redis.Z, 0)
	for _, v := range members {
		Z = append(Z, &redis.Z{
			Score:  v.Score,
			Member: v.Member,
		})
	}

	return r.Cmdable.ZAddNX(ctx, key, Z...).Result()
}

// ZCard returns the number of members in a sorted set
func (r *Client) ZCard(ctx context.Context, key string) (int64, error) {
	return r.Cmdable.ZCard(ctx, key).Result()
}

// ZCount returns the count of members in a sorted set with scores within the given values
func (r *Client) ZCount(ctx context.Context, key string, min string, max string) (int64, error) {
	return r.Cmdable.ZCount(ctx, key, min, max).Result()
}

// ZIncrBy increments the score of a member in a sorted set
func (r *Client) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	return r.Cmdable.ZIncrBy(ctx, key, increment, member).Result()
}

// ZRange returns the items of a sorted set within the range of indexes
func (r *Client) ZRange(ctx context.Context, key string, start int64, end int64) ([]string, error) {
	return r.Cmdable.ZRange(ctx, key, start, end).Result()
}

// ZRangeByScore ZRange returns the items of a sorted set within a range of scores(both inclusive)
func (r *Client) ZRangeByScore(ctx context.Context, key string, min string, max string) ([]string, error) {
	rangeBy := &redis.ZRangeBy{Min: min, Max: max}
	return r.Cmdable.ZRangeByScore(ctx, key, rangeBy).Result()
}

// ZRem removes given items from a sorted set
func (r *Client) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.Cmdable.ZRem(ctx, key, members...).Result()
}

// ZRemRangeByRank removes [start,stop] items from the given sorted set
func (r *Client) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	return r.Cmdable.ZRemRangeByRank(ctx, key, start, stop).Result()
}

// ZScore returns the score of a item in a sorted set
func (r *Client) ZScore(ctx context.Context, key string, member string) (float64, error) {
	return r.Cmdable.ZScore(ctx, key, member).Result()
}

// ZScan iterates on items of a sorted set
func (r *Client) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return r.Cmdable.ZScan(ctx, key, cursor, match, count).Result()
}
