package core

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_admin"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_auth"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_autocode"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_home"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_menu"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_role"
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"gitee.com/mxcker/go-admin/server/modules"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var Server *gin.Engine

func init() {
	Server = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	Server.Use(cors.New(corsConfig))
	Server.Use(static.Serve("/admin", static.LocalFile("./dist", true)))
}

func registerSystem(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	err := global.Init()
	if err != nil {
		return err
	}
	err = global.DB.AutoMigrate(
		sys_menu.SysMenu{},
		sys_role.SysRole{},
		sys_admin.SysAdmin{},
		sys_autocode.SysAutoCode{},
	)
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

func Start() {
	//创建API路由
	AdminAuthGroup := Server.Group("/api", middlewares.SysAuth)
	PublicGroup := Server.Group("/api")
	//注册系统组件
	err := registerSystem(AdminAuthGroup.Group("sys"), PublicGroup.Group("sys"))
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
	err = Server.Run("0.0.0.0:9000")
	if err != nil {
		panic(err)
	}
}
