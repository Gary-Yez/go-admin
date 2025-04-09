package sys_autocode

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

func Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	err := global.DB.AutoMigrate(SysAutoCode{})
	if err != nil {
		return err
	}
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
