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
	"gitee.com/mxcker/go-admin/server/modules"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"strconv"
)

var Server *gin.Engine

func registerSystemService(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	if err := initDB(); err != nil {
		return err
	}
	if err := initData(); err != nil {
		return err
	}
	sys_autocode.Register("/autocode", adminAuthGroup, publicGroup)
	sys_auth.Register("/auth", adminAuthGroup, publicGroup)
	sys_admin.Register("/admin", adminAuthGroup, publicGroup)
	sys_home.Register("/home", adminAuthGroup, publicGroup)
	sys_menu.Register("/menu", adminAuthGroup, publicGroup)
	sys_role.Register("/role", adminAuthGroup, publicGroup)
	sys_global_variable.Register("/global_variable", adminAuthGroup, publicGroup)
	return nil
}

func Start() {
	//创建服务器
	Server = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	Server.Use(cors.New(corsConfig))
	Server.Use(static.Serve("/admin", static.LocalFile("./dist", true)))
	//创建API路由
	AdminAuthGroup := Server.Group("/api", middlewares.SysAuth)
	PublicGroup := Server.Group("/api")
	//初始化全局变量
	err := global.Init()
	if err != nil {
		panic(err)
	}
	//注册系统组件
	err = registerSystemService(AdminAuthGroup.Group("sys"), PublicGroup.Group("sys"))
	if err != nil {
		panic(err)
	}
	//初始化用户组件的表
	err = modules.Init()
	if err != nil {
		panic(err)
	}
	//注册用户组件
	err = modules.Register(AdminAuthGroup, PublicGroup)
	if err != nil {
		panic(err)
	}
	// 监听并在 0.0.0.0:8080 上启动服务
	err = Server.Run(global.Config.Server.Host + ":" + strconv.Itoa(global.Config.Server.Port))
	if err != nil {
		panic(err)
	}
}
