package router

import (
	"gitee.com/mxcker/go-admin/server/core"
	"gitee.com/mxcker/go-admin/server/core/middlewares"
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

func Start() {
	//创建API路由
	AdminAuthGroup := Server.Group("/api", middlewares.SysAuth)
	PublicGroup := Server.Group("/api")
	//注册系统组件
	err := core.Register(AdminAuthGroup.Group("sys"), PublicGroup.Group("sys"))
	if err != nil {
		panic(err)
	}
	// 监听并在 0.0.0.0:8080 上启动服务
	err = Server.Run("0.0.0.0:9000")
	if err != nil {
		panic(err)
	}
}
