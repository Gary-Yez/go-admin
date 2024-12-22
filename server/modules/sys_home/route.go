package sys_home

import "github.com/gin-gonic/gin"

var controller = new(Controller)

type Route struct {
}

func (receiver *Route) Register(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) {
	Admin := adminAuthGroup.Group("/sys_home")
	{
		Admin.GET("statistic", controller.Statistic)
	}
	//无需鉴权的路由
	//Public := publicGroup.Group("/sys_home")
	//{
	//	Public.GET("statistic", controller.Statistic)
	//}
}
