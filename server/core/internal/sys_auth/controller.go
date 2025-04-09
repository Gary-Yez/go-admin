package sys_auth

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_menu"
	"gitee.com/mxcker/go-admin/server/core/types"
	"gitee.com/mxcker/go-admin/server/core/types/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"sort"
)

type controllerStruct struct {
}

func (_ *controllerStruct) Login(ctx *gin.Context) {
	body := new(loginJson)
	err := ctx.ShouldBindJSON(body)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	admin, err := Service.GetAdminByUsername(body.Username)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(body.Password))
	if err != nil {
		response.Error(ctx, "密码不正确")
	} else if admin.Status != 1 {
		response.Error(ctx, "管理员被禁用")
	} else {
		jwt := types.NewJwt(global.Config.Jwt)
		tokenString, err := jwt.Generate(admin.Id, admin.RoleId)
		if err != nil {
			response.Error(ctx, err.Error())
			return
		}
		response.Success(ctx, gin.H{
			"token": tokenString,
		})
	}
}

func (_ *controllerStruct) GetMe(ctx *gin.Context) {
	authUser, ok := ctx.MustGet("AuthUser").(types.AuthUser)
	if !ok {
		response.Error(ctx, "登录失效")
		return
	}
	user, err := Service.GetUser(authUser.UserId)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	//
	// 根据 Sort 字段对切片进行排序
	sort.Slice(user.Role.Menus, func(i, j int) bool {
		return user.Role.Menus[i].Sort < user.Role.Menus[j].Sort // 按照 Sort 字段升序排序
	})
	user.Role.Menus = sys_menu.Service.ListToTree(user.Role.Menus)
	response.Success(ctx, user)
}
