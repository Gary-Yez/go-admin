package sys_apis

import (
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"

	"github.com/gin-gonic/gin"
)

var (
	Controller = new(controllerStruct)
	Service    = new(serviceStruct)
)

type Mounter struct {
}

func (_ *Mounter) Name() string {
	return "核心服务-API管理"
}

func (_ *Mounter) Menus() []sys_menu.Definition {
	return []sys_menu.Definition{
		{
			Name: "API管理", Key: "sys_apis", ParentKey: "sys_permission", Icon: "iconoir:network",
			Path: "sys_apis", Component: "../core/views/sys_apis/index.vue", Sort: 1,
		},
	}
}

func (_ *Mounter) Initialize() error {
	// 这里执行一些初始化操作
	// 初始化数据库
	err := state.DB().AutoMigrate(&SysApi{})
	if err != nil {
		return err
	}
	err = InitData()
	if err != nil {
		return err
	}
	return nil
}

func (_ *Mounter) AdminRouter(adminGroup *gin.RouterGroup) {
	adminGroup.POST("list", Controller.List)
	adminGroup.POST("delete", Controller.Delete)
	adminGroup.POST("edit", Controller.Edit)
	adminGroup.GET("invalid_apis", Controller.GetInvalidAPIs)
	adminGroup.GET("get_groups", Controller.GetGroups)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {
	// 这里注册公共路由
}
