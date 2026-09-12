package initialization

import (
	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/config"
	"time"
)

func initCache(cfg *config.Config) (cache.Cache, error) {
	var store cache.Cache
	var err error
	if cfg.Redis.IsNotEmpty() {
		store, err = cache.NewRedisCache(cfg.Redis.Option())
	} else {
		store, err = cache.NewMemoryCache(time.Minute)
	}
	if err != nil {
		return nil, err
	}
	namespace := cache.DatabaseNamespace(cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.Database)
	return cache.WithNamespace(store, namespace), nil
}
