package global

import (
	"gitee.com/mxcker/go-admin/server/pkg/cacher"
	"gitee.com/mxcker/go-admin/server/pkg/timer"
	"gitee.com/mxcker/go-admin/server/types/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	Routes gin.RoutesInfo
	DB     *gorm.DB
	Cache  cacher.Cache
	Timer  timer.Scheduler
)
