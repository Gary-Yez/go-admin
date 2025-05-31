package test

import (
	"github.com/gin-gonic/gin"
)

var controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Name() string {
	return "测试模块"
}

func (_ *Mounter) Initialize() error {
	// 这里执行一些初始化操作
	return nil
}

func (_ *Mounter) AdminRouter(adminAuthGroup *gin.RouterGroup) {
	adminAuthGroup.GET("test", func(c *gin.Context) {})
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {

}
