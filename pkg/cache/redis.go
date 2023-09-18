package cache

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

var ErrItemNotFound = errors.New("cache: item not found")

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{
		rdb: rdb,
	}
}

func (c *RedisCache) Set(key string, value []byte, ttl time.Duration) error {
	if _, err := c.rdb.Set(key, value, ttl).Result(); err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) Get(key string) ([]byte, error) {
	return c.rdb.Get(key).Bytes()
}
