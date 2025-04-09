package sys_home

import (
	"github.com/gin-gonic/gin"
	"time"
)

var StartTime = time.Now()

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

func Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	Admin := adminAuthGroup.Group(path)
	{
		Admin.GET("statistic", Controller.Statistic)
	}
	//无需鉴权的路由
	//Public := publicGroup.Group("/sys_home")
	//{
	//	Public.GET("statistic", controller.Statistic)
	//}
	return nil
}
