package test

import (
	"github.com/gin-gonic/gin"
)

var (
	Controller = new(controllerStruct)
	Service    = new(serviceStruct)
)

type Mounter struct {
}

func (_ *Mounter) Name() string {
	return "测试模块"
}

func (_ *Mounter) Initialize() error {
	// 这里执行一些初始化操作
	return nil
}

func (_ *Mounter) AdminRouter(adminGroup *gin.RouterGroup) {
	// 这里注册需要管理权限的路由
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {
	// 这里注册公共路由
	publicGroup.GET("/test", Controller.Test)
}
