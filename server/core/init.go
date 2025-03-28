package core

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_admin"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_autocode"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_menu"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_role"
)

func Init() error {
	err := global.DB.AutoMigrate(
		sys_menu.SysMenu{},
		sys_role.SysRole{},
		sys_admin.SysAdmin{},
		sys_autocode.SysAutoCode{},
	)
	if err != nil {
		return err
	}
	sys_admin.Init()
	return nil
}
