package sys_auth

import (
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_admin"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_menu"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/pkg/request"
	"gitee.com/mxcker/go-admin/server/pkg/response"
	"gitee.com/mxcker/go-admin/server/types"
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
		jwt := types.NewJwt()
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
	authUser := request.GetAuthUser(ctx)
	user, err := Service.GetUser(authUser.UserId)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	// 根据 Sort 字段对切片进行排序
	sort.Slice(user.Role.Menus, func(i, j int) bool {
		return user.Role.Menus[i].Sort < user.Role.Menus[j].Sort // 按照 Sort 字段升序排序
	})
	user.Role.Menus = sys_menu.Service.ListToTree(user.Role.Menus)
	response.Success(ctx, user)
}

func (_ *controllerStruct) ChangeInfo(ctx *gin.Context) {
	authUser := request.GetAuthUser(ctx)
	body := struct {
		Nickname string `json:"nickname" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	user, err := Service.GetUser(authUser.UserId)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = global.DB.Model(sys_admin.SysAdmin{}).Where("id = ?", user.Id).Updates(map[string]interface{}{
		"nickname": body.Nickname,
		"phone":    body.Phone,
		"email":    body.Email,
	}).Error
	if err != nil {
		response.Error(ctx, err.Error())
		return
	} else {
		response.Success(ctx)
	}
}

func (_ *controllerStruct) ChangePassword(ctx *gin.Context) {
	authUser := request.GetAuthUser(ctx)
	body := struct {
		OldPassword     string `json:"old_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	user, err := Service.GetUser(authUser.UserId)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	if body.ConfirmPassword != body.NewPassword {
		response.Error(ctx, "两次输入的密码不一致")
		return
	}
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.OldPassword))
	if err != nil {
		response.Error(ctx, "旧密码不正确")
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	err = global.DB.Model(sys_admin.SysAdmin{}).Where("id = ?", user.Id).Updates(map[string]interface{}{
		"password_hash": password,
	}).Error
	if err != nil {
		response.Error(ctx, err.Error())
		return
	} else {
		response.Success(ctx)
	}
}
