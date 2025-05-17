package modular

import (
	"github.com/gin-gonic/gin"
)

type Mounter interface {
	Initialize() error
	Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error
}

type Loader interface {
	Add(moduleName string, mounter Mounter)
	Initialize() error
	Server(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error
}
