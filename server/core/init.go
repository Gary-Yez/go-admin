package core

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_admin"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_autocode"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_global_variable"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_menu"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_role"
)

func initCore() {

}

func initDB() error {
	err := global.DB.AutoMigrate(
		sys_menu.SysMenu{},
		sys_role.SysRole{},
		sys_admin.SysAdmin{},
		sys_autocode.SysAutoCode{},
		sys_global_variable.SysGlobalVariable{},
	)
	return err
}

func initData() error {
	if err := sys_menu.InitData(); err != nil {
		return err
	}
	if err := sys_role.InitData(); err != nil {
		return err
	}
	if err := sys_admin.InitData(); err != nil {
		return err
	}
	if err := sys_global_variable.InitData(); err != nil {
		return err
	}
	return nil
}
