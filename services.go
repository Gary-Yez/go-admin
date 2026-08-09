package admin

import (
	"gitee.com/mxcker/go-admin/cache"
	"gitee.com/mxcker/go-admin/config"
	"gitee.com/mxcker/go-admin/internal/state"
	"gitee.com/mxcker/go-admin/scheduler"
	"gorm.io/gorm"
)

func Config() *config.Config         { return state.Config() }
func DB() *gorm.DB                   { return state.DB() }
func Cache() cache.Cache             { return state.Cache() }
func Scheduler() scheduler.Scheduler { return state.Scheduler() }
