package core

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_admin"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_auth"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_autocode"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_home"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_menu"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_role"
	"github.com/gin-gonic/gin"
)

func Register(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	err := global.Init()
	if err != nil {
		return err
	}
	err = Init()
	if err != nil {
		return err
	}
	sys_autocode.Register("/autocode", adminAuthGroup, publicGroup)
	sys_auth.Register("/auth", adminAuthGroup, publicGroup)
	sys_admin.Register("/admin", adminAuthGroup, publicGroup)
	sys_home.Register("/home", adminAuthGroup, publicGroup)
	sys_menu.Register("/menu", adminAuthGroup, publicGroup)
	sys_role.Register("/role", adminAuthGroup, publicGroup)
	return nil
}
