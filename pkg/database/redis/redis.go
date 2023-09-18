package redis

import (
	"github.com/go-redis/redis"
)

func NewRedisClient(addr, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := rdb.Ping().Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}
