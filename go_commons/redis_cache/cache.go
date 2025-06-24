package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/redis"
)

type RedisCache struct {
	Client     *redis.Client
	Serializer ISerializer
	NameSpace  string
}

func NewRedisCacheClient(client *redis.Client, serializer ISerializer, nameSpace string) *RedisCache {
	if nameSpace == "" {
		panic("unable to find the service name for the redis service.")
	}

	return &RedisCache{
		Client:     client,
		Serializer: serializer,
		NameSpace:  nameSpace,
	}
}

// Set function for cache
func (r *RedisCache) Set(ctx context.Context, k string, x interface{}, d time.Duration) (done bool, err error) {
	b, err := r.Serializer.Marshal(x)
	if err != nil {
		log.Errorf("[Cache] Setting Cache key %s failed. err: %s", k, err)
		return
	}

	done, err = r.Client.Set(ctx, r.transformedKey(k), string(b), d)
	if err != nil {
		log.Errorf("[Cache] Setting Cache key %s failed. err: %s", k, err)
	}

	return
}

// Get function for cache
func (r *RedisCache) Get(ctx context.Context, k string, obj interface{}) (found bool, err error) {
	str, err := r.Client.Get(ctx, r.transformedKey(k))
	if err != nil && err != r.Client.Nil {
		log.Errorf("[Cache] Getting data from cache failed. Key %s. err: %s", k, err.Error())
		return
	}

	if err == r.Client.Nil {
		return false, nil
	}

	err = r.Serializer.Unmarshal([]byte(str), obj)
	if err != nil {
		log.Errorf("[Cache] Failed to unmarshal. Key %s. err: %s", k, err.Error())
		return false, err
	}

	return true, nil
}

func (r *RedisCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Exists(ctx, r.transformedKeys(keys)...)
}

// Delete function for cache
func (r *RedisCache) Delete(ctx context.Context, k string) (done bool, err error) {
	log.Debugf("[Cache] DEL. Key %s: ", r.transformedKey(k))

	count, err := r.Client.Del(ctx, r.transformedKey(k))
	if err != nil {
		log.Errorf("[Cache] failed to delete cache. key: %s, err: %s", k, err)
	}

	return count > 0, err
}

func (r *RedisCache) Expire(ctx context.Context, k string, expiry time.Duration) (bool, error) {
	return r.Client.Expire(ctx, r.transformedKey(k), expiry)
}

func (r *RedisCache) Pipeline() (*Pipeline, error) {
	return &Pipeline{
		pipeliner:  r.Client.Pipeline(),
		serializer: r.Serializer,
	}, nil
}

func (r *RedisCache) PipedMSet(ctx context.Context, kvArr []KVIn, d time.Duration) error {
	p, err := r.Pipeline()
	if err != nil {
		return err
	}
	defer p.Close()

	for _, kv := range kvArr {
		_, err = p.Set(ctx, r.transformedKey(kv.Key), kv.Val, d)
		if err != nil {
			log.Errorf("[Cache] PipedMSet failed. key: %s, err: %v", kv.Key, err)
			return err
		}
	}

	return p.Exec(ctx)
}

func (r *RedisCache) PipedMGet(ctx context.Context, kvArr []*KVOut) (okCount int, err error) {
	p, err := r.Pipeline()
	if err != nil {
		return 0, err
	}
	defer p.Close()

	cmds := make([]*StringCmd, len(kvArr))
	for i, kv := range kvArr {
		cmds[i], err = p.Get(ctx, r.transformedKey(kv.Key))
		if err != nil {
			log.Errorf("[Cache] PipedMGet failed. key: %s, err: %v", kv.Key, err)
			return 0, err
		}
	}

	execErr := p.Exec(ctx)
	if execErr != nil && execErr != r.Client.Nil {
		log.Errorf("[Cache] PipedExec failed. err: %v", execErr)
	}

	for i, cmd := range cmds {
		kvArr[i].exists, kvArr[i].err = cmd.Result(kvArr[i].Val)
		if kvArr[i].OK() {
			okCount++
		}
	}
	return okCount, nil
}

// MSet when we don't need to set ttl for cache keys [atomic]
func (r *RedisCache) MSet(ctx context.Context, kvArr []KVIn) error {
	var fValues []interface{}
	for _, v := range kvArr {
		fValues = append(fValues, r.transformedKey(v.Key))

		b, err := r.Serializer.Marshal(v.Val)
		if err != nil {
			log.Errorf("[Cache] MSet Serialise Cache key %v failed. err: %s", v, err)
			return err
		}
		fValues = append(fValues, string(b))
	}
	err := r.Client.MSet(ctx, fValues...)
	if err != nil {
		log.Errorf("[Cache] MSet Cache keys %v failed. err: %s", kvArr, err)
	}
	return err
}

func (r *RedisCache) ZCard(ctx context.Context, key string) (int64, error) {
	return r.Client.ZCard(ctx, r.transformedKey(key))
}

func (r *RedisCache) ZRange(ctx context.Context, key string, start int64, end int64) ([]string, error) {
	return r.Client.ZRange(ctx, r.transformedKey(key), start, end)
}

func (r *RedisCache) ZRangeByScore(ctx context.Context, key string, min string, max string) ([]string, error) {
	return r.Client.ZRangeByScore(ctx, r.transformedKey(key), min, max)
}

func (r *RedisCache) ZRem(ctx context.Context, key string, members []string) (int64, error) {
	var IMembers []interface{}
	for _, member := range members {
		IMembers = append(IMembers, member)
	}

	return r.Client.ZRem(ctx, r.transformedKey(key), IMembers...)
}

func (r *RedisCache) ZAdd(ctx context.Context, key string, items []SortedSetItem) (int64, error) {
	var ZItems []redis.ZMember
	for _, item := range items {
		ZItems = append(ZItems, redis.ZMember{
			Member: item.Member,
			Score:  item.Score,
		})
	}

	return r.Client.ZAdd(ctx, r.transformedKey(key), ZItems...)
}

func (r *RedisCache) Unlink(ctx context.Context, keys []string) (int64, error) {
	return r.Client.Unlink(ctx, r.transformedKeys(keys)...)
}

func (r *RedisCache) HSet(ctx context.Context, key string, field, value string) (int64, error) {
	return r.Client.HSet(ctx, r.transformedKey(key), field, value)
}

func (r *RedisCache) HSetAll(ctx context.Context, key string, values map[string]interface{}) (int64, error) {
	return r.Client.HSetAll(ctx, r.transformedKey(key), values)
}

func (r *RedisCache) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.Client.HGet(ctx, r.transformedKey(key), field)
}

func (r *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, r.transformedKey(key))
}

func (r *RedisCache) HDel(ctx context.Context, key string, fields []string) (int64, error) {
	return r.Client.HDel(ctx, r.transformedKey(key), fields...)
}

func (r *RedisCache) HIncrBy(ctx context.Context, key string, field string, incr int64) (int64, error) {
	return r.Client.HIncrBy(ctx, r.transformedKey(key), field, incr)
}

func (r *RedisCache) HIncrByFloat(ctx context.Context, key string, field string, incr float64) (float64, error) {
	return r.Client.HIncrByFloat(ctx, r.transformedKey(key), field, incr)
}

func (r *RedisCache) IncrBy(ctx context.Context, key string, value int64) error {
	result := r.Client.IncrBy(ctx, r.transformedKey(key), value)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *RedisCache) DecrBy(ctx context.Context, key string, value int64) error {
	result := r.Client.DecrBy(ctx, r.transformedKey(key), value)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *RedisCache) MGet(ctx context.Context, keys []string) ([]interface{}, error) {
	return r.Client.MGet(ctx, r.transformedKeys(keys)...)
}

func (r *RedisCache) IncrWithLimit(ctx context.Context, key string, limit int, expiry time.Duration) (int64, error) {
	return r.Client.IncrWithLimit(ctx, r.transformedKey(key), limit, expiry)
}

func (r *RedisCache) IncrWithMOD(ctx context.Context, key string, mod uint64, expiry time.Duration) (uint64, error) {
	return r.Client.IncrWithMOD(ctx, r.transformedKey(key), mod, expiry)
}

func (r *RedisCache) transformedKey(key string) string {
	return fmt.Sprintf("%s:%s", r.NameSpace, key)
}

func (r *RedisCache) transformedKeys(keys []string) []string {
	result := make([]string, 0)
	for _, key := range keys {
		result = append(result, fmt.Sprintf("%s:%s", r.NameSpace, key))
	}

	return result
}
