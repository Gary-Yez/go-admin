package sys_global_variable

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"github.com/gin-gonic/gin"
)

var controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Initialize() error {
	err := global.DB.AutoMigrate(SysGlobalVariable{})
	if err != nil {
		return err
	}
	err = InitData()
	if err != nil {
		return err
	}
	return nil
}

func (_ *Mounter) Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	Group := adminAuthGroup.Group(path)
	{
		Group.GET("get", controller.Get)
		Group.GET("list", controller.List)
		Group.POST("create", controller.Create)
		Group.POST("delete", controller.Delete)
		Group.POST("edit", controller.Edit)
	}
	// 无需鉴权的路由
	Public := publicGroup.Group(path)
	{
		_ = Public
	}
	return nil
}
