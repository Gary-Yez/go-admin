package modules

import (
	"gitee.com/mxcker/go-admin/server/core/global"
)

func Init() (err error) {
	err = global.DB.AutoMigrate()
	return err
}
