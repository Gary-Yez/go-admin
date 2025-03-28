package sys_autocode

import (
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"github.com/gin-gonic/gin"
)

func Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) {
	Group := adminAuthGroup.Group(path, middlewares.SysAuth)
	{
		Group.POST("generate", controller.Generate)
		Group.POST("preview", controller.Preview)
		Group.GET("history", controller.History)
		Group.GET("get_history", controller.GetHistory)
		Group.POST("delete_history", controller.DeleteHistory)
	}
}
