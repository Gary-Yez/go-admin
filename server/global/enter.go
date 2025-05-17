package global

import (
	"gitee.com/mxcker/go-admin/server/pkg/cacher"
	"gitee.com/mxcker/go-admin/server/pkg/timer"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Cache cacher.Cache
	Timer timer.Scheduler
)
