package sys_admin

import (
	"gitee.com/mxcker/go-admin/internal/state"

	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Name() string {
	return "核心服务-管理员"
}

func (_ *Mounter) AdminRouter(adminAuthGroup *gin.RouterGroup) {
	adminAuthGroup.GET("get", Controller.Get)
	adminAuthGroup.GET("list", Controller.List)
	adminAuthGroup.POST("create", Controller.Create)
	adminAuthGroup.POST("delete", Controller.Delete)
	adminAuthGroup.POST("edit", Controller.Edit)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {

}

func (_ *Mounter) Initialize() error {
	err := state.DB().AutoMigrate(&SysAdmin{})
	if err != nil {
		return err
	}
	if err = InitData(); err != nil {
		return err
	}
	return nil
}
