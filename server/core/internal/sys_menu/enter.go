package sys_menu

import "github.com/gin-gonic/gin"

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

func Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) {
	Group := adminAuthGroup.Group(path)
	{
		Group.GET("get", Controller.Get)
		Group.GET("list", Controller.List)
		Group.POST("create", Controller.Create)
		Group.POST("delete", Controller.Delete)
		Group.POST("edit", Controller.Edit)
	}
}
