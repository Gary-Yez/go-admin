package sys_admin

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

func Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	err := global.DB.AutoMigrate(&SysAdmin{})
	if err != nil {
		return err
	}
	if err = InitData(); err != nil {
		return err
	}
	Group := adminAuthGroup.Group(path)
	{
		Group.GET("get", Controller.Get)
		Group.GET("list", Controller.List)
		Group.POST("create", Controller.Create)
		Group.POST("delete", Controller.Delete)
		Group.POST("edit", Controller.Edit)
	}
	return nil
}
