package sys_apis

import (
	"gitee.com/mxcker/go-admin/internal/state"

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

func (_ *Mounter) Initialize() error {
	// 这里执行一些初始化操作
	// 初始化数据库
	err := state.DB().AutoMigrate(&SysApi{}, &SysIgnoreApi{})
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
	adminGroup.GET("get", Controller.Get)
	adminGroup.POST("list", Controller.List)
	adminGroup.POST("create", Controller.Create)
	adminGroup.POST("delete", Controller.Delete)
	adminGroup.POST("edit", Controller.Edit)
	adminGroup.POST("update_ignore", Controller.UpdateIgnore)
	adminGroup.GET("sync_api", Controller.SyncApi)
	adminGroup.GET("get_groups", Controller.GetGroups)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {
	// 这里注册公共路由
}
