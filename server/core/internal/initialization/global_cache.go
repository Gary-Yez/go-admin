package initialization

import (
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/pkg/cacher"
	"github.com/redis/go-redis/v9"
	"time"
)

func initCache() error {
	if global.Config.Redis.IsNotEmpty() {
		r, err := cacher.NewRedisCache(&redis.Options{
			Addr:     global.Config.Redis.Host + ":" + global.Config.Redis.Port,
			Username: global.Config.Redis.Username,
			Password: global.Config.Redis.Password, // 没有密码默认留空
			DB:       global.Config.Redis.DB,       // 默认使用0号数据库
		})
		if err != nil {
			return err
		}
		global.Cache = r
		return nil
	} else {
		c, err := cacher.NewMemoryCache(time.Minute)
		if err != nil {
			return err
		}
		global.Cache = c
		return nil
	}
}
