package modular

import (
	"github.com/gin-gonic/gin"
)

type Mounter interface {
	Initialize() error
	AdminRouter(adminGroup *gin.RouterGroup)
	PublicRouter(publicGroup *gin.RouterGroup)
}

type Loader interface {
	Mount(routerPrefix string, mounter Mounter)
	InitializeAll() error
	RegisterRouter(adminGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error
}
