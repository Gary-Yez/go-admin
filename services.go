package admin

import (
	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/scheduler"
	"github.com/Gary-Yez/go-admin/internal/state"
	"gorm.io/gorm"
)

// CacheStore 和 CacheLock 提供业务可复用的缓存接口，具体实现由框架管理。
type CacheStore = cache.Cache
type CacheLock = cache.Lock

var ErrCacheNotFound = cache.ErrCacheNotFound

// 定时任务注册需要的类型通过 admin 公开。
type HandlerOption = scheduler.HandlerOption
type HandlerParams = scheduler.HandlerParams
type HandlerParam = scheduler.HandlerParam
type ParamsType = scheduler.ParamsType

const (
	StringParams = scheduler.StringParams
	IntParams    = scheduler.IntParams
	BoolParams   = scheduler.BoolParams
)

func DB() *gorm.DB                   { return state.DB() }
func Cache() cache.Cache             { return state.Cache() }
func Scheduler() scheduler.Scheduler { return state.Scheduler() }
