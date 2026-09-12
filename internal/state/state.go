package state

import (
	"sync"

	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/config"
	"github.com/Gary-Yez/go-admin/internal/scheduler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var current struct {
	sync.RWMutex
	config      *config.Config
	db          *gorm.DB
	cache       cache.Cache
	scheduler   scheduler.Scheduler
	adminRoutes gin.RoutesInfo
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

func SetAdminRoutes(routes gin.RoutesInfo) {
	current.Lock()
	defer current.Unlock()
	current.adminRoutes = append(gin.RoutesInfo(nil), routes...)
}

func AdminRoutes() gin.RoutesInfo {
	current.RLock()
	defer current.RUnlock()
	return append(gin.RoutesInfo(nil), current.adminRoutes...)
}
