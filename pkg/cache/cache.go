package cache

import "time"

type Cache interface {
	Set(key string, value []byte, ttl time.Duration) error
	Get(key string) ([]byte, error)
}
