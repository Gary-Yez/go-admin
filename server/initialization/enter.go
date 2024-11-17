package initialization

import (
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/models"
)

func Init() {
	err := global.DB.AutoMigrate(
		// 系统自带
		&models.SysMenu{},
		&models.SysRole{},
		&models.SysAdmin{},
		&models.SysAutoCode{},
		// 用户自定义
	)
	if err != nil {
		panic(err)
	}
	InitSysMenu()
	InitSysRole()
	InitSysAdmin()
}
