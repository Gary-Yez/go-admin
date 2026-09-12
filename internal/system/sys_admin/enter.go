package sys_admin

import (
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"

	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Name() string {
	return "核心服务-管理员"
}

func (_ *Mounter) Menus() []sys_menu.Definition {
	return []sys_menu.Definition{
		{
			Name: "管理员管理", Key: "sys_admin", ParentKey: "sys_permission", Icon: "iconoir:group",
			Path: "sys_admin", Component: "../core/views/sys_admin/index.vue", Sort: 3,
		},
	}
}

func (_ *Mounter) AdminRouter(adminAuthGroup *gin.RouterGroup) {
	adminAuthGroup.GET("list", Controller.List)
	adminAuthGroup.POST("create", Controller.Create)
	adminAuthGroup.POST("delete", Controller.Delete)
	adminAuthGroup.POST("edit", Controller.Edit)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {

}

func (_ *Mounter) Initialize() error {
	if err := state.DB().SetupJoinTable(&SysAdmin{}, "Roles", &SysAdminRole{}); err != nil {
		return err
	}
	err := state.DB().AutoMigrate(&SysAdmin{})
	if err != nil {
		return err
	}
	return InitData()
}
