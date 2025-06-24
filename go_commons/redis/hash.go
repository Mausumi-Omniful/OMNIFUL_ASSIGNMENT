package redis

import "context"

// HSet Sets a single field of a hash map to the provided value
// Returns the number of fields that were newly added
func (r *Client) HSet(ctx context.Context, key string, field string, value interface{}) (int64, error) {
	return r.Cmdable.HSet(ctx, key, map[string]interface{}{field: value}).Result()
}

// HSetAll Sets a multiple fields of a hash map to the provided value
// Returns the number of fields that were newly added
func (r *Client) HSetAll(ctx context.Context, key string, values map[string]interface{}) (int64, error) {
	return r.Cmdable.HSet(ctx, key, values).Result()
}

// HGet Gets a particular field value in a hash map
func (r *Client) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.Cmdable.HGet(ctx, key, field).Result()
}

// HGetAll Gets the entire hash map corresponding to a key
func (r *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.Cmdable.HGetAll(ctx, key).Result()
}

// HMGet Gets the values from a particular key's multiple fields
func (r *Client) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	return r.Cmdable.HMGet(ctx, key, fields...).Result()
}

// HDel Delete the values from a particular key's multiple fields
func (r *Client) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return r.Cmdable.HDel(ctx, key, fields...).Result()
}

func (r *Client) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	return r.Cmdable.HIncrBy(ctx, key, field, incr).Result()
}

func (r *Client) HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error) {
	return r.Cmdable.HIncrByFloat(ctx, key, field, incr).Result()
}
