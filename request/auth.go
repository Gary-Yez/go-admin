package request

import (
	"gitee.com/mxcker/go-admin/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetAuthUser(ctx *gin.Context) *utils.AuthUser {
	user, has := ctx.Get("AuthUser")
	if has {
		authUser, ok := user.(*utils.AuthUser)
		if ok {
			return authUser
		} else {
			return &utils.AuthUser{}
		}
	} else {
		return &utils.AuthUser{}
	}
}

func SetAuthUser(ctx *gin.Context, authUser *utils.AuthUser) {
	ctx.Set("AuthUser", authUser)
}
