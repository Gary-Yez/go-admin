package bizconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Gary-Yez/go-admin/internal/state"
	"gorm.io/gorm"
)

// 先列出 Key，每项拿到写锁后重新读数据库，避免将查询期间的旧值写回缓存。
func SyncCache(ctx context.Context) (int, error) {
	var keys []string
	if err := state.DB().WithContext(ctx).Model(&SysConfigValue{}).Order("config_key").Pluck("config_key", &keys).Error; err != nil {
		return 0, err
	}
	for i, key := range keys {
		if err := syncCachedValue(ctx, key); err != nil {
			return i, fmt.Errorf("已同步 %d 项，配置 %s 同步失败：%w", i, key, err)
		}
	}
	return len(keys), nil
}

func syncCachedValue(ctx context.Context, key string) error {
	lock, err := lockValue(ctx, key)
	if err != nil {
		return err
	}
	defer lock.Unlock()
	var row SysConfigValue
	err = state.DB().WithContext(ctx).First(&row, "config_key = ?", key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return state.Cache().Del(valueCacheKey(key))
	}
	if err != nil {
		return err
	}
	value := json.RawMessage(row.Value)
	if err := ValidateValue(row.Type, value); err != nil {
		return err
	}
	return cacheValue(key, value)
}

func CleanupInvalid(ctx context.Context, keys []string) (int, error) {
	if registeredKeys == nil {
		return 0, errors.New("配置结构尚未初始化，不能清理")
	}
	// 删除前先检查全部请求项，不允许通过该接口删除当前实例使用的配置。
	for _, key := range keys {
		if _, exists := registeredKeys[key]; exists {
			return 0, fmt.Errorf("配置 %s 仍在当前实例注册，不能删除", key)
		}
	}
	count := 0
	seen := map[string]bool{}
	for _, key := range keys {
		if seen[key] {
			continue
		}
		seen[key] = true
		deleted, err := deleteInvalid(ctx, key)
		if err != nil {
			return count, fmt.Errorf("已清理 %d 项，配置 %s 清理失败：%w", count, key, err)
		}
		if deleted {
			count++
		}
	}
	return count, nil
}

func deleteInvalid(ctx context.Context, key string) (bool, error) {
	lock, err := lockValue(ctx, key)
	if err != nil {
		return false, err
	}
	defer lock.Unlock()
	// 先清缓存；失败时保留数据库数据。读取回填使用同一把锁。
	if err := state.Cache().Del(valueCacheKey(key)); err != nil {
		return false, err
	}
	result := state.DB().WithContext(ctx).Where("config_key = ?", key).Delete(&SysConfigValue{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
