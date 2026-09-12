package sys_auth

import (
	"github.com/Gary-Yez/go-admin/response"
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
	adminAuthGroup.GET("password_policy", Controller.PasswordPolicy)
	adminAuthGroup.Use(func(ctx *gin.Context) {
		if ctx.GetBool("api_token_auth") {
			response.Error(ctx, "请使用登录会话管理个人账号", 403)
		}
	})
	adminAuthGroup.POST("switch_role", Controller.SwitchRole)
	adminAuthGroup.POST("change_info", Controller.ChangeInfo)
	adminAuthGroup.POST("change_password", Controller.ChangePassword)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {
	publicGroup.POST("login", Controller.Login)
}
