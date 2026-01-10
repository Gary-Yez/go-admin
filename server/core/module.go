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

type registeredModule struct {
	Key string
	Module
}

var moduleKeyMap = make(map[string]bool)
var Modules []*registeredModule

func AddModule(moduleKey string, m Module) {
	if moduleKeyMap[moduleKey] {
		panic("重复注册的模组：" + moduleKey)
	}
	Modules = append(Modules, &registeredModule{
		Key:    moduleKey,
		Module: m,
	})
}
