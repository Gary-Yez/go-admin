package cache

import (
	"context"
	"time"
)

// DefaultTTL 是配置项、用户身份和 API 密钥身份缓存的默认有效期。
const DefaultTTL = 30 * time.Minute

type Lock interface {
	Unlock() error
}

type Cache interface {
	// Client 返回原始客户端，直接调用时不会自动添加缓存命名空间
	Client() interface{}
	// Set 写入键值对，ttl为0则缓存永不过期
	Set(key string, value interface{}, ttl time.Duration) error
	// SetJSON 把结构体序列化并写入缓存，ttl为0则永不过期
	SetJSON(key string, data any, ttl time.Duration) error
	// Get 字符串读取
	Get(key string) (value string, err error)
	// MustGet Get的忽略错误版本
	MustGet(key string) string
	// GetBool Bool读取
	GetBool(key string) (value bool, err error)
	// MustGetBool GetBool的忽略错误版本
	MustGetBool(key string) bool
	// GetInt Int读取
	GetInt(key string) (value int, err error)
	// MustGetInt GetInt的忽略错误版本
	MustGetInt(key string) int
	// GetJSON 把缓存中的数据反序列化成结构体
	GetJSON(key string, dest any) error
	// Exists 检查 key 是否存在
	Exists(key string) (bool, error)
	// Del 删除 key
	Del(key string) error
	// Lock 只尝试一次获取锁，ctx 控制本次获取操作：
	// expiration 锁自动释放时间
	// extendInterval 锁的自动续期间隔
	Lock(ctx context.Context, key string, expiration time.Duration, extendInterval time.Duration) (Lock, error)
	// WaitForLock 等待并重试获取锁，直到成功或 ctx 取消、超时。
	WaitForLock(ctx context.Context, key string, expiration time.Duration, extendInterval time.Duration) (Lock, error)
}

// waitForLock 为内存和 Redis 共用的等待逻辑。
func waitForLock(ctx context.Context, store Cache, key string, expiration, extendInterval time.Duration) (Lock, error) {
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		lock, err := store.Lock(ctx, key, expiration, extendInterval)
		if err == nil {
			if err := ctx.Err(); err != nil {
				_ = lock.Unlock()
				return nil, err
			}
			return lock, nil
		}
		timer := time.NewTimer(30 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}
