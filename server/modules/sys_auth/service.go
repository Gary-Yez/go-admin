package sys_auth

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/models"
	"gorm.io/gorm"
)

type Service struct {
}

func (_ *Service) GetAdminByUsername(username string) (admin *models.SysAdmin, err error) {
	admin = &models.SysAdmin{}
	err = global.DB.Model(models.SysAdmin{}).Where("username = ?", username).First(admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("管理员不存在")
	}
	return admin, err
}

func (_ *Service) GetUser(id uint) (user *models.SysAdmin, err error) {
	user = &models.SysAdmin{}
	err = global.DB.Preload("Role").Where("id = ?", id).First(user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("管理员账号不存在")
	}
	if user.Status != 1 {
		return nil, errors.New("管理员被禁用")
	}
	var menus []*models.SysMenu
	if user.Role != nil && user.Role.Default {
		if err = global.DB.Model(&models.SysMenu{}).Find(&menus).Error; err != nil {
			return nil, err
		}
		user.Role.Menus = menus
	} else {
		return user, global.DB.Model(&models.SysRole{}).Preload("Menus").Find(user.Role).Error
	}
	return user, err
}
