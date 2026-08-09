package admin

import (
	"github.com/Gary-Yez/go-admin/cache"
	"github.com/Gary-Yez/go-admin/config"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/scheduler"
	"gorm.io/gorm"
)

func Config() *config.Config         { return state.Config() }
func DB() *gorm.DB                   { return state.DB() }
func Cache() cache.Cache             { return state.Cache() }
func Scheduler() scheduler.Scheduler { return state.Scheduler() }
