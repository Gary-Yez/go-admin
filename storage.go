package admin

import (
	"fmt"
	"github.com/Gary-Yez/go-admin/internal/system/sys_storage"
	"github.com/Gary-Yez/go-admin/storage"
)

// Storage 按账号 ID 创建适配器，不传参数时使用已启用的默认存储。
// 本地适配器使用完毕后通过 io.Closer 关闭。
func Storage(ids ...uint) (storage.Store, error) {
	config, err := StorageConfig(ids...)
	if err != nil {
		return nil, err
	}
	return storage.New(config)
}

// StorageConfig 从存储管理读取账号配置，优先缓存。返回值包含凭据，不要返回前端。
func StorageConfig(ids ...uint) (storage.Config, error) {
	if len(ids) > 1 {
		return storage.Config{}, fmt.Errorf("只能指定一个存储账号")
	}
	id := uint(0)
	if len(ids) == 1 {
		id = ids[0]
	}
	return sys_storage.ReadConfig(id, true)
}
