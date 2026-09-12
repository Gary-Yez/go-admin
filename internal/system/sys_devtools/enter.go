package sys_devtools

import (
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"

	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct{}

func (m *Mounter) Name() string {
	return "核心服务-开发工具"
}

func (_ *Mounter) Initialize() error {
	err := state.DB().AutoMigrate(SysAutoCode{})
	if err != nil {
		return err
	}
	return nil
}

func (_ *Mounter) AdminRouter(_ *gin.RouterGroup) {}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {
	if state.Config() == nil || !state.Config().IsDev() {
		return
	}
	publicGroup.POST("generate", Controller.Generate)
	publicGroup.GET("menu_options", sys_menu.Controller.List)
	publicGroup.POST("preview", Controller.Preview)
	publicGroup.GET("history", Controller.History)
	publicGroup.GET("get_history", Controller.GetHistory)
	publicGroup.POST("delete_history", Controller.DeleteHistory)
	publicGroup.POST("preview_delete_history", Controller.PreviewDeleteHistory)
}
