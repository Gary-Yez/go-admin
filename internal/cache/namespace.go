package cache

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"
)

// DatabaseNamespace 为缓存、锁和通知频道生成相同的数据库命名空间。
// MD5 仅用于生成数据库标识，不用于安全认证或完整性校验。
func DatabaseNamespace(host, port, database string) string {
	identity := fmt.Sprintf("%s:%s/%s", host, port, database)
	return fmt.Sprintf("go-admin:%x:", md5.Sum([]byte(identity)))
}

// WithNamespace 为缓存和锁统一添加前缀，业务只传模块与具体键。
func WithNamespace(store Cache, namespace string) Cache {
	return &namespacedCache{store: store, prefix: namespace}
}

type namespacedCache struct {
	store  Cache
	prefix string
}

func (c *namespacedCache) cacheKey(key string) string { return c.prefix + "cache:" + key }

// Client 返回原始客户端；直接操作客户端不会自动添加命名空间。
func (c *namespacedCache) Client() interface{} { return c.store.Client() }

func (c *namespacedCache) Set(key string, value interface{}, ttl time.Duration) error {
	return c.store.Set(c.cacheKey(key), value, ttl)
}

func (c *namespacedCache) SetJSON(key string, data any, ttl time.Duration) error {
	return c.store.SetJSON(c.cacheKey(key), data, ttl)
}

func (c *namespacedCache) Get(key string) (string, error) { return c.store.Get(c.cacheKey(key)) }

func (c *namespacedCache) MustGet(key string) string { return c.store.MustGet(c.cacheKey(key)) }

func (c *namespacedCache) GetBool(key string) (bool, error) { return c.store.GetBool(c.cacheKey(key)) }

func (c *namespacedCache) MustGetBool(key string) bool { return c.store.MustGetBool(c.cacheKey(key)) }

func (c *namespacedCache) GetInt(key string) (int, error) { return c.store.GetInt(c.cacheKey(key)) }

func (c *namespacedCache) MustGetInt(key string) int { return c.store.MustGetInt(c.cacheKey(key)) }

func (c *namespacedCache) GetJSON(key string, dest any) error {
	return c.store.GetJSON(c.cacheKey(key), dest)
}

func (c *namespacedCache) Exists(key string) (bool, error) { return c.store.Exists(c.cacheKey(key)) }

func (c *namespacedCache) Del(key string) error { return c.store.Del(c.cacheKey(key)) }

func (c *namespacedCache) Lock(ctx context.Context, key string, expiration, extendInterval time.Duration) (Lock, error) {
	return c.store.Lock(ctx, c.prefix+"lock:"+key, expiration, extendInterval)
}

func (c *namespacedCache) WaitForLock(ctx context.Context, key string, expiration, extendInterval time.Duration) (Lock, error) {
	return c.store.WaitForLock(ctx, c.prefix+"lock:"+key, expiration, extendInterval)
}
