package sys_menu

import (
	"github.com/Gary-Yez/go-admin/internal/state"

	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Name() string {
	return "核心服务-菜单管理"
}

func (_ *Mounter) Menus() []Definition {
	return []Definition{
		{
			Name: "仪表盘", Key: "sys_home", Icon: "iconoir:dashboard",
			Path: "sys_home", Component: "../core/views/sys_home/index.vue", Sort: 0,
		},
		{
			Name: "权限管理", Key: "sys_permission", Icon: "iconoir:shield-check",
			Path: "sys_permission", Sort: 1,
		},
		{
			Name: "系统运维", Key: "sys_operations", Icon: "iconoir:server", Path: "sys_operations", Sort: 2,
		},
		{
			Name: "菜单管理", Key: "sys_menu", ParentKey: "sys_permission", Icon: "iconoir:tree",
			Path: "sys_menu", Component: "../core/views/sys_menu/index.vue", Sort: 0,
		},
	}
}

func (_ *Mounter) Initialize() error {
	return state.DB().AutoMigrate(&SysMenu{})
}

func (_ *Mounter) AdminRouter(adminAuthGroup *gin.RouterGroup) {
	adminAuthGroup.GET("list", Controller.List)
	adminAuthGroup.POST("create", Controller.Create)
	adminAuthGroup.POST("delete", Controller.Delete)
	adminAuthGroup.POST("edit", Controller.Edit)
	adminAuthGroup.POST("sort", Controller.Sort)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {
}
