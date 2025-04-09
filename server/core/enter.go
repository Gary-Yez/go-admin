package core

import (
	"gitee.com/mxcker/go-admin/server/core/global"
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
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Start(Server *gin.Engine) {
	Server.Use(static.Serve("/admin", static.LocalFile("./dist", true)))
	//创建API路由
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
	// 监听并在 0.0.0.0:8080 上启动服务
	err = Server.Run(global.Config.Server.Host + ":" + strconv.Itoa(global.Config.Server.Port))
	if err != nil {
		panic(err)
	}
}
