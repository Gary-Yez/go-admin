package sys_auth

import (
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Initialize() error {
	return nil
}
func (_ *Mounter) Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	Group := adminAuthGroup.Group(path)
	{
		Group.GET("me", middlewares.SysAuth, Controller.GetMe)
		Group.POST("change_info", middlewares.SysAuth, Controller.ChangeInfo)
		Group.POST("change_password", middlewares.SysAuth, Controller.ChangePassword)
	}
	Public := publicGroup.Group(path)
	{
		Public.POST("login", Controller.Login)
	}
	return nil
}
