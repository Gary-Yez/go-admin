package bizconfig

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/state"
)

func valueCacheKey(key string) string {
	return "config:" + key
}

// 写入和未命中回填共用 cache.Lock，避免旧值在写入完成后覆盖缓存。
func lockValue(ctx context.Context, key string) (cache.Lock, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	lock, err := state.Cache().WaitForLock(ctx, valueCacheKey(key), 30*time.Second, 10*time.Second)
	if err != nil {
		return nil, errors.New("配置正被更新或缓存不可用，请稍后重试")
	}
	return lock, nil
}

func databaseValue(ctx context.Context, key string) (json.RawMessage, error) {
	var row SysConfigValue
	if err := state.DB().WithContext(ctx).First(&row, "config_key = ?", key).Error; err != nil {
		return nil, err
	}
	value := json.RawMessage(row.Value)
	if err := ValidateValue(row.Type, value); err != nil {
		return nil, err
	}
	return value, nil
}

func cacheValue(key string, value json.RawMessage) error {
	if err := state.Cache().Set(valueCacheKey(key), string(value), cache.DefaultTTL); err != nil {
		// 数据库提交后缓存写入失败，尽量移除旧值，下次读取重新回填。
		return errors.Join(err, state.Cache().Del(valueCacheKey(key)))
	}
	return nil
}

func ReadValue(key string) (json.RawMessage, error) {
	store := state.Cache()
	if store == nil || state.DB() == nil {
		return nil, errors.New("系统配置尚未初始化")
	}
	cached, cacheErr := store.Get(valueCacheKey(key))
	if cacheErr == nil && cached != "null" && json.Valid([]byte(cached)) {
		return json.RawMessage(cached), nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	if cacheErr != nil && !errors.Is(cacheErr, cache.ErrCacheNotFound) {
		return databaseValue(ctx, key)
	}
	lock, err := lockValue(ctx, key)
	if err != nil {
		return databaseValue(ctx, key)
	}
	defer lock.Unlock()
	// 等待锁期间可能已有请求回填，避免重复访问数据库。
	if cached, err := store.Get(valueCacheKey(key)); err == nil && cached != "null" && json.Valid([]byte(cached)) {
		return json.RawMessage(cached), nil
	}
	value, err := databaseValue(ctx, key)
	if err != nil {
		return nil, err
	}
	// 已读取有效数据库值时，回填失败不影响本次读取。
	_ = cacheValue(key, value)
	return value, nil
}
