package request

import (
	"gitee.com/mxcker/go-admin/server/types"
	"github.com/gin-gonic/gin"
)

func GetAuthUser(ctx *gin.Context) *types.AuthUser {
	user, has := ctx.Get("AuthUser")
	if has {
		authUser, ok := user.(*types.AuthUser)
		if ok {
			return authUser
		} else {
			return &types.AuthUser{}
		}
	} else {
		return &types.AuthUser{}
	}
}

func SetAuthUser(ctx *gin.Context, authUser *types.AuthUser) {
	ctx.Set("AuthUser", authUser)
}
