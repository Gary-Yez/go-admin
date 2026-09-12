package middlewares

import (
	"github.com/Gary-Yez/go-admin/internal/permissions"
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
		authUser, err := request.GetAuthUser(c)
		if err != nil {
			response.Error(c, err.Error(), 401)
			return
		}
		// 获取请求路径和方法
		path := strings.TrimPrefix(c.Request.URL.Path, apiPrefix)
		method := c.Request.Method
		// 自身鉴权接口只免除角色授权，账号状态与角色归属仍由上游校验。
		if permissions.IsAuthenticatedAPI(method, path) {
			c.Next()
			return
		}
		// 检查权限
		ok, _ := enforcer.Enforce(strconv.Itoa(int(authUser.RoleId)), path, method)
		if !ok {
			response.Error(c, "无权限访问当前接口", 403)
			return
		}
		c.Next()
	}
}
