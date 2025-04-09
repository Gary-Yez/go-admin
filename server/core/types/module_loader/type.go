package module_loader

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type Mounter interface {
	Initialize() error
	Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error
}

type Loader interface {
	Add(moduleName string, mounter Mounter)
	InitAll() error
	RegisterAll(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error
}

func NewLoader() Loader {
	return &defaultLoader{
		sequence: []string{},
		mp:       make(map[string]Mounter),
		lock:     sync.Mutex{},
	}
}
