package sys_storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/storage"
)

type cachedStorage struct {
	Config  storage.Config
	Enabled bool
}

func cacheKey(id uint) string { return fmt.Sprintf("storage:account:%d", id) }
func accountLock(ctx context.Context, id uint) (cache.Lock, error) {
	wait, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return state.Cache().WaitForLock(wait, cacheKey(id), 30*time.Second, 10*time.Second)
}
func databaseConfig(ctx context.Context, id uint) (cachedStorage, error) {
	var row SysStorage
	if err := state.DB().WithContext(ctx).First(&row, id).Error; err != nil {
		return cachedStorage{}, err
	}
	config, err := row.Config()
	return cachedStorage{Config: config, Enabled: row.Enabled}, err
}

// ReadConfig 缓存按存储 ID 保存整份配置。已有文件忽略 Enabled，停用只阻止新上传。
func ReadConfig(id uint, requireEnabled bool) (storage.Config, error) {
	if state.DB() == nil || state.Cache() == nil {
		return storage.Config{}, errors.New("存储管理尚未初始化")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	if id == 0 {
		if !requireEnabled {
			return storage.Config{}, errors.New("文件没有关联存储账号")
		}
		var row SysStorage
		if err := state.DB().WithContext(ctx).Select("id").Where("default_slot = ? AND enabled = ?", 1, true).First(&row).Error; err != nil {
			return storage.Config{}, errors.New("请先设置已启用的默认存储")
		}
		id = row.Id
	}

	value, err := readCached(ctx, id)
	if err != nil {
		return storage.Config{}, err
	}
	if requireEnabled && !value.Enabled {
		return storage.Config{}, errors.New("该存储已停用")
	}
	return value.Config, nil
}

// readCached 与账号写入共用未命中回填锁，缓存异常时直接读取数据库。
func readCached(ctx context.Context, id uint) (cachedStorage, error) {
	var value cachedStorage
	err := state.Cache().GetJSON(cacheKey(id), &value)
	if err == nil {
		return value, nil
	}
	if !errors.Is(err, cache.ErrCacheNotFound) {
		return databaseConfig(ctx, id)
	}
	lock, err := accountLock(ctx, id)
	if err != nil {
		return databaseConfig(ctx, id)
	}
	defer lock.Unlock()
	if err := state.Cache().GetJSON(cacheKey(id), &value); err == nil {
		return value, nil
	}
	value, err = databaseConfig(ctx, id)
	if err == nil {
		_ = state.Cache().SetJSON(cacheKey(id), value, cache.DefaultTTL)
	}
	return value, err
}
