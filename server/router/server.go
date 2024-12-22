package router

import (
	"gitee.com/mxcker/go-admin/server/middlewares"
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
	AdminAuthGroup := Server.Group("/api", middlewares.SysAuth)
	PublicGroup := Server.Group("/api")
	for _, route := range routes {
		route.Register(AdminAuthGroup, PublicGroup)
	}
	// 监听并在 0.0.0.0:8080 上启动服务
	err := Server.Run("0.0.0.0:9000")
	if err != nil {
		return
	}
}
