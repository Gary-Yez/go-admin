package global

import (
	"gitee.com/mxcker/go-admin/server/config"
	"gitee.com/mxcker/go-admin/server/core/pkg/cacher"
	"gitee.com/mxcker/go-admin/server/core/pkg/timer"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
)

var (
	RootPath string
	Config   *config.Config
	Routes   gin.RoutesInfo
	DB       *gorm.DB
	Cache    cacher.Cache
	Timer    timer.Scheduler
)

func init() {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	RootPath = rootPath
}
