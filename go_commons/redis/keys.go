package redis

import "context"

// Del deletes the key from redis
// Notice: Unlink is always a better choice compared to Del
func (r *Client) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.Cmdable.Del(ctx, keys...).Result()
}

// Unlink lazily deletes keys from redis. It does not block the main thread like Del
// so is always a better choice compared to Del
func (r *Client) Unlink(ctx context.Context, keys ...string) (int64, error) {
	return r.Cmdable.Unlink(ctx, keys...).Result()
}

// Exists checks for existence of each key supplied and returns the number of
// keys that exist
//
// Note: If all the keys do not belong to the same slot, it will return error
func (r *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.Cmdable.Exists(ctx, keys...).Result()
}

func (r *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.Cmdable.Keys(ctx, pattern).Result()
}
