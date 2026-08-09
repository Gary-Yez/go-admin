package middlewares

import (
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func CasbinMiddleware(apiPrefix string, enforcer *casbin.SyncedCachedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色（根据你的认证系统实现）
		authUser := request.GetAuthUser(c)
		// 获取请求路径和方法
		path := strings.TrimPrefix(c.Request.URL.Path, apiPrefix)
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
