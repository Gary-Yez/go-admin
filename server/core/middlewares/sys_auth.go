package middlewares

import (
	"gitee.com/mxcker/go-admin/server/pkg/response"
	"gitee.com/mxcker/go-admin/server/types"
	"github.com/gin-gonic/gin"
	"strings"
)

func SysAuth(ctx *gin.Context) {
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		accessToken, _ := ctx.GetQuery("access_token")
		authorization = "Bearer " + accessToken
	}
	if strings.HasPrefix(authorization, "Bearer ") {
		token := strings.Split(authorization, " ")[1]
		jwt := types.NewJwt()
		accessToken, err := jwt.Parse(token)
		if err != nil {
			response.Error(ctx, err, 401)
			return
		}
		if accessToken.UserId == 0 {
			response.Error(ctx, "用户不存在", 401)
			return
		}
		ctx.Set("AuthUser", accessToken.AuthUser)
	} else {
		response.Error(ctx, "令牌格式应为Bearer开头", 401)
	}
}
