package sys_auth

import (
	"gitee.com/mxcker/go-admin/server/middlewares"
	"github.com/gin-gonic/gin"
)

var controller = new(Controller)

type Route struct {
}

func (receiver *Route) Register(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) {
	Group := adminAuthGroup.Group("/sys_auth")
	{
		Group.GET("me", middlewares.SysAuth, controller.GetMe)
	}
	Public := publicGroup.Group("/sys_auth")
	{
		Public.POST("login", controller.Login)
	}

}
