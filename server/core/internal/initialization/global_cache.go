package initialization

import (
	"gitee.com/mxcker/go-admin/server/configs"
	"gitee.com/mxcker/go-admin/server/global"
	cacher2 "gitee.com/mxcker/go-admin/server/pkg/cacher"
	"github.com/redis/go-redis/v9"
	"time"
)

func initCache() error {
	if configs.Config.Redis.Host != "" && configs.Config.Redis.Port != "" {
		r, err := cacher2.NewRedisCache(&redis.Options{
			Addr:     configs.Config.Redis.Host + ":" + configs.Config.Redis.Port,
			Username: configs.Config.Redis.Username,
			Password: configs.Config.Redis.Password, // 没有密码默认留空
			DB:       configs.Config.Redis.DB,       // 默认使用0号数据库
		})
		if err != nil {
			return err
		}
		global.Cache = r
		return nil
	} else {
		c, err := cacher2.NewMemoryCache(time.Minute)
		if err != nil {
			return err
		}
		global.Cache = c
		return nil
	}
}
