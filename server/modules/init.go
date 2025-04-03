package modules

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/modules/test"
)

func Init() (err error) {
	err = global.DB.AutoMigrate(
		&test.Test{},
		&test.Test{},
		&test.Test{},
	)
	return err
}
