package redis

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Advanced Counter Operations

// INCRWithLimit increments a key with a limit in Redis.
// It provides atomic increment operation with a maximum limit and expiry.
//
// Parameters:
//   - ctx: Context for the operation
//   - key: Redis key to increment
//   - limit: Maximum value allowed for the counter
//   - expiry: Time duration after which the key will expire
//
// Returns:
//   - newValue: Current value after increment (or current value if limit reached)
//   - err: Error if operation fails or limit is reached
//
// The operation is atomic and handles the following cases:
//  1. Key doesn't exist: Creates key with value 1
//  2. Current value >= limit: Returns error
//  3. Normal case: Increments value
func (r *Client) IncrWithLimit(
	ctx context.Context,
	key string,
	limit int,
	expiry time.Duration,
) (newValue int64, err error) {
	// Lua script for atomic increment with limit check
	luaScript := `
        local current = redis.call('GET', KEYS[1])
        if not current then
            redis.call('SET', KEYS[1], 1)
            redis.call('EXPIRE', KEYS[1], tonumber(ARGV[1]))
            return 1
        elseif tonumber(current) >= tonumber(ARGV[2]) then
            return "unable to increment"
        else
            return redis.call('INCR', KEYS[1])
        end
    `

	cmd := r.Cmdable.Eval(ctx, luaScript, []string{key}, expiry.Seconds(), limit)
	result, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	// Handle different result types
	switch result := result.(type) {
	case int64:
		return result, nil
	case string:
		return 0, errors.New(result)
	default:
		return 0, fmt.Errorf("unexpected result type: %T", result)
	}
}

// IncrWithMOD increments a key with modulo operation in Redis.
// It provides atomic increment operation with modulo, useful for circular counters.
//
// Parameters:
//   - ctx: Context for the operation
//   - key: Redis key to increment
//   - mod: Modulo value (must be > 0)
//   - expiry: Time duration after which the key will expire (0 means no expiry)
//
// Returns:
//   - newValue: (current value + 1) % mod
//   - err: Error if operation fails or if mod is 0
//
// The operation is atomic and handles the following cases:
//  1. Key doesn't exist: Creates key with value 1 % mod
//  2. Normal case: Increments value and applies modulo while preserving TTL
func (r *Client) IncrWithMOD(
	ctx context.Context,
	key string,
	mod uint64,
	expiry time.Duration,
) (newValue uint64, err error) {
	if mod == 0 {
		return 0, errors.New("mod must be greater than 0")
	}

	// Lua script for atomic increment with modulo
	luaScript := `
        local current = redis.call('GET', KEYS[1])
        local ttl = redis.call('TTL', KEYS[1])
        if not current then
            redis.call('SET', KEYS[1], 1 % tonumber(ARGV[1]))
            if tonumber(ARGV[2]) ~= 0 then
                redis.call('EXPIRE', KEYS[1], tonumber(ARGV[2]))
            end
            return 1 % tonumber(ARGV[1])
        else
            redis.call('SET', KEYS[1], (current + 1) % tonumber(ARGV[1]))
            if ttl > 0 then
                redis.call('EXPIRE', KEYS[1], ttl)
            end
            return ((current + 1) % tonumber(ARGV[1]))
        end
    `

	cmd := r.Cmdable.Eval(ctx, luaScript, []string{key}, mod, expiry.Seconds())
	rawResult, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	// Handle different result types
	switch result := rawResult.(type) {
	case int64:
		if result < 0 {
			return 0, fmt.Errorf("unexpected negative result: %d", result)
		}
		return uint64(result), nil
	case string:
		return 0, errors.New(result)
	default:
		return 0, fmt.Errorf("unexpected result type: %T", result)
	}
}
