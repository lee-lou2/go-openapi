package config

import (
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	cacheInstance *redis.Client
	cacheOnce     sync.Once
)

// GetCache 캐시 조회
func GetCache() *redis.Client {
	cacheOnce.Do(func() {
		cacheInstance = redis.NewClient(&redis.Options{
			Addr:     GetEnv("REDIS_ADDR"),
			Password: GetEnv("REDIS_PASSWORD"),
			DB:       0,
		})
	})
	return cacheInstance
}
