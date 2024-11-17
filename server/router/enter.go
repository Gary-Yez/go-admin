package router

import (
	"gitee.com/mxcker/go-admin/server/modules/sys_admin"
	"gitee.com/mxcker/go-admin/server/modules/sys_auth"
	"gitee.com/mxcker/go-admin/server/modules/sys_autocode"
	"gitee.com/mxcker/go-admin/server/modules/sys_menu"
	"gitee.com/mxcker/go-admin/server/modules/sys_role"
	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup)
}

var routes = []Route{
	// 系统自带
	new(sys_admin.Route),
	new(sys_auth.Route),
	new(sys_autocode.Route),
	new(sys_menu.Route),
	new(sys_role.Route),
	// 用户自定义
}
