package config

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

// Cache 인터페이스 정의
type Cache interface {
	Set(key string, value interface{}, expiration time.Duration)
	Get(key string) (interface{}, bool)
}

type GoCache struct {
	c *cache.Cache
}

// Set 데이터 저장
func (g *GoCache) Set(key string, value interface{}, expiration time.Duration) {
	g.c.Set(key, value, expiration*time.Second)
}

// Get 데이터 조회
func (g *GoCache) Get(key string) (interface{}, bool) {
	return g.c.Get(key)
}

// Delete 데이터 삭제
func (g *GoCache) Delete(key string) {
	g.c.Delete(key)
}

var (
	cacheInstance *GoCache
	cacheOnce     sync.Once
)

// GetCache 캐시 조회
func GetCache() *GoCache {
	cacheOnce.Do(func() {
		cacheInstance = &GoCache{
			c: cache.New(5*time.Minute, 10*time.Minute),
		}
	})
	return cacheInstance
}
