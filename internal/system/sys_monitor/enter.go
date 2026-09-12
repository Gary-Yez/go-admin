package sys_monitor

import (
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct{}

func (*Mounter) Name() string      { return "系统运维-节点监控" }
func (*Mounter) Initialize() error { return nil }
func (*Mounter) Menus() []sys_menu.Definition {
	return []sys_menu.Definition{{Name: "节点监控", Key: "sys_monitor", ParentKey: "sys_operations", Icon: "iconoir:server", Path: "sys_monitor", Component: "../core/views/sys_monitor/index.vue", Sort: 0}}
}
func (*Mounter) AdminRouter(group *gin.RouterGroup) { group.GET("list", Controller.List) }
func (*Mounter) PublicRouter(*gin.RouterGroup)      {}
