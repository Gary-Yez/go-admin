package sys_auth

import (
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

func Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	Group := adminAuthGroup.Group(path)
	{
		Group.GET("me", middlewares.SysAuth, Controller.GetMe)
	}
	Public := publicGroup.Group(path)
	{
		Public.POST("login", Controller.Login)
	}
	return nil
}
