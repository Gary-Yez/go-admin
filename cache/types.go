package cache

import "time"

type Lock interface {
	Unlock() error
}

type Cache interface {
	// Client 返回核心的客户端
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
	// Lock 分布式锁：
	// expiration 锁自动释放时间
	// extendInterval 锁的自动续期间隔
	Lock(key string, expiration time.Duration, extendInterval time.Duration) (Lock, error)
}
