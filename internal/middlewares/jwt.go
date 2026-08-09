package middlewares

import (
	"github.com/Gary-Yez/go-admin/internal/utils"
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
	"strings"
)

type APITokenVerifier func(string) (*utils.AuthUser, error)

func JWTMiddleware(secret string, verifyAPIToken APITokenVerifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		//if authorization == "" {
		//	accessToken, _ := ctx.GetQuery("access_token")
		//	authorization = "Bearer " + accessToken
		//}
		if strings.HasPrefix(authorization, "Bearer ") {
			token := strings.Split(authorization, " ")[1]
			if strings.HasPrefix(token, "API_") {
				authUser, err := verifyAPIToken(token)
				if err != nil {
					response.Error(ctx, err, 401)
					return
				}
				request.SetAuthUser(ctx, authUser)
			} else {
				jwt := utils.NewJwt(secret)
				accessToken, err := jwt.Parse(token)
				if err != nil {
					response.Error(ctx, err, 401)
					return
				}
				if accessToken.UserId == 0 {
					response.Error(ctx, "用户不存在", 401)
					return
				}
				request.SetAuthUser(ctx, &accessToken.AuthUser)
			}
		} else {
			response.Error(ctx, "令牌格式应为Bearer开头", 401)
		}
	}
}
