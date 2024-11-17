package sys_menu

import "github.com/gin-gonic/gin"

var controller = new(Controller)

type Route struct {
}

func (receiver *Route) Register(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) {
	Group := adminAuthGroup.Group("/sys_menu")
	{
		Group.GET("get", controller.Get)
		Group.GET("list", controller.List)
		Group.POST("create", controller.Create)
		Group.POST("delete", controller.Delete)
		Group.POST("edit", controller.Edit)
	}
}
