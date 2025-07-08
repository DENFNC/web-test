package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisCache{rdb: rdb}
}

func (c *RedisCache) Get(key string) ([]byte, bool) {
	ctx := context.Background()
	val, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return val, true
}

func (c *RedisCache) Set(key string, value []byte, ttl time.Duration) {
	ctx := context.Background()
	c.rdb.Set(ctx, key, value, ttl)
}

func (c *RedisCache) Invalidate(keys ...string) {
	ctx := context.Background()
	if len(keys) == 0 {
		c.rdb.FlushDB(ctx)
		return
	}
	c.rdb.Del(ctx, keys...)
}

func (c *RedisCache) Clear() {
	ctx := context.Background()
	c.rdb.FlushDB(ctx)
}

func (c *RedisCache) Delete(key string) {
	ctx := context.Background()
	c.rdb.Del(ctx, key)
}
