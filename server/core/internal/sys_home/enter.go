package sys_home

import "github.com/gin-gonic/gin"

func Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) {
	Admin := adminAuthGroup.Group(path)
	{
		Admin.GET("statistic", controller.Statistic)
	}
	//无需鉴权的路由
	//Public := publicGroup.Group("/sys_home")
	//{
	//	Public.GET("statistic", controller.Statistic)
	//}
}
