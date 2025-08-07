package sys_auth

import (
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Name() string {
	return "核心服务-鉴权服务"
}

func (_ *Mounter) Initialize() error {
	return nil
}

func (_ *Mounter) AdminRouter(adminAuthGroup *gin.RouterGroup) {
	adminAuthGroup.GET("me", Controller.GetMe)
	adminAuthGroup.POST("change_info", Controller.ChangeInfo)
	adminAuthGroup.POST("change_password", Controller.ChangePassword)
	adminAuthGroup.POST("reset_api_token", Controller.ResetApiToken)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {
	publicGroup.POST("login", Controller.Login)
}
