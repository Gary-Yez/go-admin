package core

import (
	"github.com/gin-gonic/gin"
)

type Module interface {
	Name() string
	Initialize() error
	AdminRouter(adminGroup *gin.RouterGroup)
	PublicRouter(publicGroup *gin.RouterGroup)
}

var Modules = make(map[string]Module)

func AddModule(routePath string, m Module) {
	if Modules[routePath] != nil {
		panic("重复注册的routePath：" + routePath)
	}
	Modules[routePath] = m
}
