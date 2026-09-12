package initialization

import (
	"errors"
	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/config"
	"github.com/Gary-Yez/go-admin/internal/scheduler"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB        *gorm.DB
	Cache     cache.Cache
	Scheduler scheduler.Scheduler
}

func InitDependencies(cfg *config.Config) (*Dependencies, error) {
	db, err := initGormMysql(cfg)
	if err != nil {
		return nil, errors.New("数据库初始化失败：" + err.Error())
	}
	cacheStore, err := initCache(cfg)
	if err != nil {
		return nil, errors.New("缓存初始化失败：" + err.Error())
	}
	scheduler, err := initTaskManager(cacheStore)
	if err != nil {
		return nil, errors.New("任务管理器初始化失败：" + err.Error())
	}
	return &Dependencies{DB: db, Cache: cacheStore, Scheduler: scheduler}, nil
}
