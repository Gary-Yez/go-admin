package core

import (
	"gitee.com/mxcker/go-admin/server/core/internal/sys_admin"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_auth"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_autocode"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_global_variable"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_home"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_menu"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_role"
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"gitee.com/mxcker/go-admin/server/core/models"
	"gitee.com/mxcker/go-admin/server/modules"
	"github.com/gin-gonic/gin"
)

func Load(Server *gin.Engine) {
	AdminAuthGroup := Server.Group("/api", middlewares.SysAuth)
	PublicGroup := Server.Group("/api")
	//注册系统组件
	moduleMap := new(models.ModuleMap)
	moduleMap.Add("sys/autocode", sys_autocode.Register)
	moduleMap.Add("sys/global_variable", sys_global_variable.Register)
	moduleMap.Add("sys/menu", sys_menu.Register)
	moduleMap.Add("sys/role", sys_role.Register)
	moduleMap.Add("sys/admin", sys_admin.Register)
	moduleMap.Add("sys/auth", sys_auth.Register)
	moduleMap.Add("sys/home", sys_home.Register)
	//加载用户组件
	err := modules.Load(moduleMap)
	if err != nil {
		panic(err)
	}
	err = moduleMap.RegisterAll(AdminAuthGroup, PublicGroup)
	if err != nil {
		panic(err)
	}
}
