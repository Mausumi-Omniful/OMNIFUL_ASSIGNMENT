package redis

import (
	"context"
	"time"
)

// Set sets the value against a given redis key
// Zero expiration means the key has no expiration time
func (r *Client) Set(ctx context.Context, key string, val string, expiration time.Duration) (bool, error) {
	result, err := r.Cmdable.Set(ctx, key, val, expiration).Result()
	if err == nil && result == "OK" {
		return true, nil
	}

	return false, err
}

// Get returns the data from the redis key
func (r *Client) Get(ctx context.Context, k string) (string, error) {
	return r.Cmdable.Get(ctx, k).Result()
}

// MGet returns the data corresponding to the list of redis keys
func (r *Client) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return r.Cmdable.MGet(ctx, keys...).Result()
}

// MSet sets multiple k, v pairs into redis
// It throws an error if keys are not corresponding to the same slot
func (r *Client) MSet(ctx context.Context, pairs ...interface{}) error {
	return r.Cmdable.MSet(ctx, pairs...).Err()
}

// Incr Increments the number stored at redis key by one
func (r *Client) Incr(ctx context.Context, k string) (int64, error) {
	return r.Cmdable.Incr(ctx, k).Result()
}

// Decr Decrements the number stored at redis key by one
func (r *Client) Decr(ctx context.Context, k string) (int64, error) {
	return r.Cmdable.Decr(ctx, k).Result()
}

// Expire Set explicit expiry for a key
func (r *Client) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.Cmdable.Expire(ctx, key, expiration).Result()
}

// SetNX Sets a key only when the key doesn't exist
func (r *Client) SetNX(ctx context.Context, key string, val string, expiration time.Duration) (bool, error) {
	return r.Cmdable.SetNX(ctx, key, val, expiration).Result()
}

// TTL returns the time left for expiry from the redis key
func (r *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.Cmdable.TTL(ctx, key).Result()
}
