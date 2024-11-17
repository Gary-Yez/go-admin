package sys_auth

import (
	"gitee.com/mxcker/go-admin/server/models"
	"gitee.com/mxcker/go-admin/server/modules/sys_menu"
	"gitee.com/mxcker/go-admin/server/utils"
	"gitee.com/mxcker/go-admin/server/utils/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"sort"
)

var service = new(Service)
var menuService = new(sys_menu.Service)

type Controller struct {
}

func (_ *Controller) Login(ctx *gin.Context) {
	body := new(loginJson)
	err := ctx.ShouldBindJSON(body)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	admin, err := service.GetAdminByUsername(body.Username)
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
		jwt := utils.NewJwt()
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

func (_ *Controller) GetMe(ctx *gin.Context) {
	authUser, ok := ctx.MustGet("AuthUser").(models.AuthUser)
	if !ok {
		response.Error(ctx, "登录失效")
		return
	}
	user, err := service.GetUser(authUser.UserId)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	// 根据 Sort 字段对切片进行排序
	sort.Slice(user.Role.Menus, func(i, j int) bool {
		return user.Role.Menus[i].Sort < user.Role.Menus[j].Sort // 按照 Sort 字段升序排序
	})
	user.Role.Menus = menuService.ListToTree(user.Role.Menus)
	response.Success(ctx, user)
}
