package sys_role

import (
	"github.com/gin-gonic/gin"
)

var controller = new(controllerStruct)
var Service = new(serviceStruck)

func Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) {
	Group := adminAuthGroup.Group(path)
	{
		Group.GET("get", controller.Get)
		Group.GET("list", controller.List)
		Group.POST("create", controller.Create)
		Group.POST("delete", controller.Delete)
		Group.POST("edit", controller.Edit)
		Group.POST("permission", controller.UpdatePermission)
	}
}
