package global

import (
	"gitee.com/mxcker/go-admin/server/core/configs"
	"gorm.io/gorm"
)

var (
	RootPath string
	Config   *configs.Config
	DB       *gorm.DB
	Redis    string
)
