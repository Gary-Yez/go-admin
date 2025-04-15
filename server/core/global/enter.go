package global

import (
	"gitee.com/mxcker/go-admin/server/core/pkg/cacher"
	"gitee.com/mxcker/go-admin/server/core/pkg/timer"
	"gorm.io/gorm"
	"os"
)

var (
	DB    *gorm.DB
	Cache cacher.Cache
	Timer timer.Scheduler
)

func IsDev() bool {
	//return true
	return os.Getenv("APP_ENV") != "production"
}
