package middlewares

import (
	"context"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/utils"
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
	"strings"
)

type APITokenVerifier func(context.Context, string) (*utils.AuthUser, error)
type AuthUserVerifier func(*utils.AuthUser) error

func JWTMiddleware(verifyAPIToken APITokenVerifier, verifyAuthUser AuthUserVerifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		if strings.HasPrefix(authorization, "Bearer ") {
			token := strings.Split(authorization, " ")[1]
			if strings.HasPrefix(token, "API_") {
				authUser, err := verifyAPIToken(ctx.Request.Context(), token)
				if err != nil {
					response.Error(ctx, err, 401)
					return
				}
				ctx.Set("api_token_auth", true)
				request.SetAuthUser(ctx, authUser)
			} else {
				accessToken, err := utils.NewJwt(state.Config().JWT.Secret).Parse(token)
				if err != nil {
					response.Error(ctx, err, 401)
					return
				}

				request.SetAuthUser(ctx, &accessToken.AuthUser)
			}
		} else {
			response.Error(ctx, "令牌格式应为Bearer开头", 401)
			return
		}
		authUser, err := request.GetAuthUser(ctx)
		if err != nil {
			response.Error(ctx, err.Error(), 401)
			return
		}
		// JWT 和 API 密钥统一核对账号状态与角色归属。
		if err := verifyAuthUser(authUser); err != nil {
			response.Error(ctx, err.Error(), 401)
			return
		}
	}
}
