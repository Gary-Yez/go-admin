package request

import (
	"errors"
	"github.com/Gary-Yez/go-admin/internal/utils"
	"github.com/gin-gonic/gin"
)

var ErrUnauthenticated = errors.New("登录身份缺失或无效")

func GetAuthUser(ctx *gin.Context) (*utils.AuthUser, error) {
	if ctx == nil {
		return nil, ErrUnauthenticated
	}
	user, has := ctx.Get("AuthUser")
	if has {
		authUser, ok := user.(*utils.AuthUser)
		if ok && authUser != nil && authUser.UserId > 0 && authUser.RoleId > 0 {
			return authUser, nil
		}
	}
	return nil, ErrUnauthenticated
}

func SetAuthUser(ctx *gin.Context, authUser *utils.AuthUser) {
	ctx.Set("AuthUser", authUser)
}
