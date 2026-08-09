package sys_devtools

import (
	"gitee.com/mxcker/go-admin/internal/state"

	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
	name      string
	routePath string
}

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

func (_ *Mounter) AdminRouter(adminAuthGroup *gin.RouterGroup) {
	adminAuthGroup.POST("generate", Controller.Generate)
	adminAuthGroup.POST("preview", Controller.Preview)
	adminAuthGroup.GET("history", Controller.History)
	adminAuthGroup.POST("delete_history", Controller.DeleteHistory)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {

}
