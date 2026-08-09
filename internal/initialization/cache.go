package initialization

import (
	cacher2 "github.com/Gary-Yez/go-admin/cache"
	"github.com/Gary-Yez/go-admin/config"
	"github.com/redis/go-redis/v9"
	"time"
)

func initCache(cfg *config.Config) (cacher2.Cache, error) {
	if cfg.Redis.IsNotEmpty() {
		r, err := cacher2.NewRedisCache(&redis.Options{
			Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
			Username: cfg.Redis.Username,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		if err != nil {
			return nil, err
		}
		return r, nil
	} else {
		c, err := cacher2.NewMemoryCache(time.Minute)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}
