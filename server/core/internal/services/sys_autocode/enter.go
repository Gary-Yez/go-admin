package sys_autocode

import (
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"gitee.com/mxcker/go-admin/server/global"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Initialize() error {
	err := global.DB.AutoMigrate(SysAutoCode{})
	if err != nil {
		return err
	}
	return nil
}

func (_ *Mounter) Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {

	Group := adminAuthGroup.Group(path, middlewares.SysAuth)
	{
		Group.POST("generate", Controller.Generate)
		Group.POST("preview", Controller.Preview)
		Group.GET("history", Controller.History)
		Group.GET("get_history", Controller.GetHistory)
		Group.POST("delete_history", Controller.DeleteHistory)
	}
	return nil
}
