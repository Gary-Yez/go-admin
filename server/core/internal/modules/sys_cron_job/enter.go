package sys_cron_job

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"github.com/gin-gonic/gin"
)

var controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Initialize() error {
	// 初始化数据库
	err := global.DB.AutoMigrate(&SysCronJob{})
	if err != nil {
		return err
	}
	return nil
}

func (_ *Mounter) Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	// 注册路由
	// 需要管理员权限的路由
	Group := adminAuthGroup.Group(path)
	{
		//Group.GET("get_registered_task", controller.GetRegisteredTask)
		Group.GET("get", controller.Get)
		Group.POST("list", controller.List)
		Group.POST("create", controller.Create)
		Group.POST("delete", controller.Delete)
		Group.POST("edit", controller.Edit)
	}
	// 无需鉴权的路由
	Public := publicGroup.Group(path)
	{
		_ = Public
		Public.GET("get_registered_handler", controller.GetRegisteredHandler)

	}
	return nil
}
