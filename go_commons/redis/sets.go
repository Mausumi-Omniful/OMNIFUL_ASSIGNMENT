package redis

import "context"

// SAdd takes a list of elements and add them in the set
// Returns number of elements added
func (r *Client) SAdd(ctx context.Context, key string, elements ...interface{}) (int64, error) {
	return r.Cmdable.SAdd(ctx, key, elements...).Result()
}

// SIsMember checks whether element exists in set or not
// Returns TRUE/FALSE is element exists in set or not
func (r *Client) SIsMember(ctx context.Context, key string, element interface{}) (bool, error) {
	return r.Cmdable.SIsMember(ctx, key, element).Result()
}
