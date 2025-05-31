package sys_devtools

import (
	"gitee.com/mxcker/go-admin/server/global"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Name() string {
	return "核心服务-自动代码生成"
}

func (_ *Mounter) Initialize() error {
	err := global.DB.AutoMigrate(SysAutoCode{})
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
