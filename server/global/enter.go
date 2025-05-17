package global

import (
	"gitee.com/mxcker/go-admin/server/core/pkg/cacher"
	"gitee.com/mxcker/go-admin/server/core/pkg/timer"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Cache cacher.Cache
	Timer timer.Scheduler
)
