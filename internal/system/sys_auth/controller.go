package sys_auth

import (
	"errors"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_admin"
	"github.com/Gary-Yez/go-admin/internal/system/sys_login_log"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/Gary-Yez/go-admin/internal/utils"
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"sort"
	"time"
)

type controllerStruct struct {
}

func (_ *controllerStruct) PasswordPolicy(ctx *gin.Context) {
	minimum, err := utils.ReadLoginInt("login.password_min_length")
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	ctx.Header("Cache-Control", "no-store")
	response.Success(ctx, gin.H{"min_length": minimum, "max_bytes": utils.PasswordMaxBytes})
}

func (_ *controllerStruct) Login(ctx *gin.Context) {
	started := time.Now()
	body := new(loginJson)
	entry := &sys_login_log.SysLoginLog{CreatedAt: started, IP: ctx.RemoteIP(), UserAgent: ctx.Request.UserAgent(), Status: sys_login_log.Failed, Message: "登录请求异常"}
	defer func() {
		entry.Username = body.Username
		entry.Duration = time.Since(started).Milliseconds()
		if err := sys_login_log.Service.Record(entry); err != nil {
			log.Printf("写入登录日志失败：%v", err)
		}
	}()
	if err := ctx.ShouldBindJSON(body); err != nil {
		entry.Message = "登录参数不正确"
		response.Error(ctx, err.Error())
		return
	}
	// 当前未配置可信代理，使用连接来源，防止伪造转发头绕过限制。
	user, err := Service.Login(ctx.Request.Context(), body.Username, body.Password, entry.IP)
	if err != nil {
		entry.Message = err.Error()
		if errors.Is(err, errLoginLocked) {
			entry.Status = sys_login_log.Blocked
		}
		response.Error(ctx, err.Error())
		return
	}
	token, err := generateLoginToken(user.Id, user.RoleId, user.LoginVersion)
	if err != nil {
		entry.Message = "登录令牌签发失败"
		response.Error(ctx, err.Error())
		return
	}
	entry.UserId = user.Id
	entry.Status = sys_login_log.Success
	entry.Message = "登录成功"
	response.Success(ctx, gin.H{"token": token})
}

func (_ *controllerStruct) GetMe(ctx *gin.Context) {
	authUser, err := request.GetAuthUser(ctx)
	if err != nil {
		response.Error(ctx, err.Error(), 401)
		return
	}
	user, err := Service.GetUser(authUser.UserId, authUser.RoleId)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	// 根据 Sort 字段对切片进行排序
	sort.Slice(user.Role.Menus, func(i, j int) bool {
		if user.Role.Menus[i].Sort == user.Role.Menus[j].Sort {
			return user.Role.Menus[i].Id < user.Role.Menus[j].Id
		}
		return user.Role.Menus[i].Sort < user.Role.Menus[j].Sort // 按照 Sort 字段升序排序
	})
	user.Role.Menus = sys_menu.Service.ListToTree(user.Role.Menus)
	response.Success(ctx, user)
}

func (_ *controllerStruct) SwitchRole(ctx *gin.Context) {
	authUser, err := request.GetAuthUser(ctx)
	if err != nil {
		response.Error(ctx, err, 401)
		return
	}
	body := struct {
		RoleId uint `json:"role_id" binding:"required"`
	}{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := Service.VerifyAuthUser(&utils.AuthUser{UserId: authUser.UserId, RoleId: body.RoleId}); err != nil {
		response.Error(ctx, "不可切换到未绑定的角色", 403)
		return
	}
	token, err := generateLoginToken(authUser.UserId, body.RoleId, *authUser.LoginVersion)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"token": token})
}

func (_ *controllerStruct) ChangeInfo(ctx *gin.Context) {
	authUser, err := request.GetAuthUser(ctx)
	if err != nil {
		response.Error(ctx, err.Error(), 401)
		return
	}
	body := struct {
		Nickname string `json:"nickname" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = state.DB().Model(sys_admin.SysAdmin{}).Where("id = ?", authUser.UserId).Updates(map[string]interface{}{
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
	authUser, err := request.GetAuthUser(ctx)
	if err != nil {
		response.Error(ctx, err.Error(), 401)
		return
	}
	body := struct {
		OldPassword     string `json:"old_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	var user sys_admin.SysAdmin
	if err := state.DB().Select("id", "password_hash", "login_version").First(&user, authUser.UserId).Error; err != nil {
		response.Error(ctx, err.Error())
		return
	}
	if body.ConfirmPassword != body.NewPassword {
		response.Error(ctx, "两次输入的密码不一致")
		return
	}
	if err := utils.ValidatePassword(body.NewPassword); err != nil {
		response.Error(ctx, err.Error())
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
	// 先签发，避免签名配置读取失败时密码已经修改。
	token, err := generateLoginToken(user.Id, authUser.RoleId, user.LoginVersion+1)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if authUser.LoginVersion == nil || *authUser.LoginVersion != user.LoginVersion {
		response.Error(ctx, "登录已失效，请重新登录", 401)
		return
	}
	if err := sys_admin.Service.ChangePassword(user.Id, user.PasswordHash, string(password), user.LoginVersion); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"token": token})
}
