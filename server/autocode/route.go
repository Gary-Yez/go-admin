package autocode

import (
	"gitee.com/mxcker/go-admin/server/middlewares"
	"github.com/gin-gonic/gin"
)

var controller = new(Controller)

type Route struct {
}

func (receiver *Route) Register(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) {
	Group := adminAuthGroup.Group("/autocode", middlewares.SysAuth)
	{
		Group.POST("generate", controller.Generate)
		Group.POST("preview", controller.Preview)
		Group.GET("history", controller.History)
		Group.GET("get_history", controller.GetHistory)
		Group.POST("delete_history", controller.DeleteHistory)
	}
}
