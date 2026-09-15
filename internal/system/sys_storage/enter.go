package sys_storage

import (
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct{}

func (*Mounter) Name() string      { return "文件管理-存储管理" }
func (*Mounter) Initialize() error { return state.DB().AutoMigrate(&SysStorage{}) }
func (*Mounter) Menus() []sys_menu.Definition {
	return []sys_menu.Definition{{Name: "存储管理", Key: "sys_storage", ParentKey: "sys_files", Icon: "iconoir:database", Path: "sys_storage", Component: "../core/views/sys_storage/index.vue", Sort: 1}}
}
func (*Mounter) PublicRouter(*gin.RouterGroup) {}
func (*Mounter) AdminRouter(group *gin.RouterGroup) {
	group.POST("list", Controller.List)
	group.GET("get", Controller.Get)
	group.POST("save", Controller.Save)
	group.POST("default", Controller.SetDefault)
	group.POST("enabled", Controller.SetEnabled)
	group.POST("delete", Controller.Delete)
}
