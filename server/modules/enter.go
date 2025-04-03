package modules

import (
	"github.com/gin-gonic/gin"
	"gitee.com/mxcker/go-admin/server/modules/test"
)

func Register(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	test.Register("test", adminAuthGroup, publicGroup)
	test.Register("test", adminAuthGroup, publicGroup)
	test.Register("test", adminAuthGroup, publicGroup)
	return nil
}
