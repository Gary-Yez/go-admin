package state

import (
	"sync"

	"gitee.com/mxcker/go-admin/cache"
	"gitee.com/mxcker/go-admin/config"
	"gitee.com/mxcker/go-admin/scheduler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var current struct {
	sync.RWMutex
	config    *config.Config
	db        *gorm.DB
	cache     cache.Cache
	scheduler scheduler.Scheduler
	routes    gin.RoutesInfo
}

func Configure(cfg *config.Config, db *gorm.DB, store cache.Cache, jobs scheduler.Scheduler) {
	current.Lock()
	defer current.Unlock()
	current.config = cfg
	current.db = db
	current.cache = store
	current.scheduler = jobs
}

func Config() *config.Config { current.RLock(); defer current.RUnlock(); return current.config }
func DB() *gorm.DB           { current.RLock(); defer current.RUnlock(); return current.db }
func Cache() cache.Cache     { current.RLock(); defer current.RUnlock(); return current.cache }
func Scheduler() scheduler.Scheduler {
	current.RLock()
	defer current.RUnlock()
	return current.scheduler
}

func SetRoutes(routes gin.RoutesInfo) {
	current.Lock()
	defer current.Unlock()
	current.routes = routes
}

func Routes() gin.RoutesInfo {
	current.RLock()
	defer current.RUnlock()
	return append(gin.RoutesInfo(nil), current.routes...)
}
