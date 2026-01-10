package middlewares

import (
	"gitee.com/mxcker/go-admin/server/core/pkg/request"
	"gitee.com/mxcker/go-admin/server/core/pkg/response"
	"gitee.com/mxcker/go-admin/server/global"
	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func CasbinMiddleware(enforcer *casbin.SyncedCachedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色（根据你的认证系统实现）
		authUser := request.GetAuthUser(c)
		// 获取请求路径和方法
		path := strings.TrimPrefix(c.Request.URL.Path, global.Config.Server.ApiPrefix)
		method := c.Request.Method
		// 检查权限
		ok, _ := enforcer.Enforce(strconv.Itoa(int(authUser.RoleId)), path, method)
		if !ok {
			response.Error(c, "无权限访问当前接口", 403)
			return
		}
		c.Next()
	}
}
