package sys_config

import (
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct{}

func (*Mounter) Name() string      { return "系统运维-配置管理" }
func (*Mounter) Initialize() error { return nil }
func (*Mounter) Menus() []sys_menu.Definition {
	return []sys_menu.Definition{{Name: "配置管理", Key: "sys_config", ParentKey: "sys_operations", Icon: "iconoir:settings", Path: "sys_config", Component: "../core/views/sys_config/values.vue", Sort: 3}}
}
func (*Mounter) AdminRouter(group *gin.RouterGroup) {
	group.GET("values", Controller.Values)
	group.POST("update_value", Controller.UpdateValue)
	group.POST("reset_value", Controller.ResetValue)
	group.POST("sync_cache", Controller.SyncCache)
	group.POST("cleanup_invalid", Controller.CleanupInvalid)
}
func (*Mounter) PublicRouter(group *gin.RouterGroup) {
	group.GET("site", Controller.Site)
	if state.Config() == nil || !state.Config().IsDev() {
		return
	}
	group.GET("list", Controller.List)
	group.POST("preview", Controller.Preview)
	group.POST("apply", Controller.Apply)
}
