package sys_auth

import (
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"github.com/gin-gonic/gin"
)

func Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) {
	Group := adminAuthGroup.Group(path)
	{
		Group.GET("me", middlewares.SysAuth, controller.GetMe)
	}
	Public := publicGroup.Group(path)
	{
		Public.POST("login", controller.Login)
	}

}
