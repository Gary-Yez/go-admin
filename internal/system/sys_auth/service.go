package sys_auth

import (
	"errors"
	"github.com/Gary-Yez/go-admin/internal/state"
	utils2 "github.com/Gary-Yez/go-admin/internal/utils"

	"github.com/Gary-Yez/go-admin/internal/system/sys_admin"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/Gary-Yez/go-admin/internal/system/sys_role"
	"gorm.io/gorm"
	"time"
)

type serviceStruct struct {
}

func (s *serviceStruct) GetAdminByUsername(username string) (admin *sys_admin.SysAdmin, err error) {
	admin = &sys_admin.SysAdmin{}
	err = state.DB().Model(sys_admin.SysAdmin{}).Where("username = ?", username).First(admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("管理员不存在")
	}
	return admin, err
}

func (s *serviceStruct) GetUser(id uint) (user *sys_admin.SysAdmin, err error) {
	user = &sys_admin.SysAdmin{}
	err = state.DB().Preload("Role").Where("id = ?", id).First(user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("管理员账号不存在")
	}
	if user.Status != 1 {
		return nil, errors.New("管理员被禁用")
	}
	var menus []*sys_menu.SysMenu
	if user.Role != nil && user.Role.Default {
		if err = state.DB().Model(&sys_menu.SysMenu{}).Find(&menus).Error; err != nil {
			return nil, err
		}
		user.Role.Menus = menus
	} else {
		return user, state.DB().Model(&sys_role.SysRole{}).Preload("Menus").Find(user.Role).Error
	}
	return user, err
}

func (s *serviceStruct) ResetApiToken(id uint) (string, error) {
	randomString, err := utils2.GenerateRandomString(32)
	if err != nil {
		return "", err
	}
	randomString = "API_" + randomString
	return randomString, state.DB().Model(sys_admin.SysAdmin{}).Where("id = ?", id).UpdateColumn("api_token", randomString).Error
}

func (s *serviceStruct) VerifyApiToken(token string) (user *utils2.AuthUser, err error) {
	user = &utils2.AuthUser{}
	if token == "" {
		return nil, errors.New("无效的API密钥")
	}
	cacheKey := "core:api_token_cache:" + token
	// 如果命中缓存则取返回值
	err = state.Cache().GetJSON(cacheKey, user)
	if err == nil && user.UserId != 0 && user.RoleId != 0 {
		return user, nil
	}
	// 未命中则查询数据库并且写入缓存
	admin := &sys_admin.SysAdmin{}
	err = state.DB().Model(&sys_admin.SysAdmin{}).Where("api_token = ?", token).First(admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("无效的APIToken")
		}
		return nil, errors.New("API验证失败")
	}
	authUser := utils2.AuthUser{
		UserId: admin.Id,
		RoleId: admin.RoleId,
	}
	_ = state.Cache().SetJSON(cacheKey, authUser, time.Hour)
	return &authUser, nil
}
