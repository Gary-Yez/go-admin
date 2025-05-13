package sys_auth

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/internal/modules/sys_admin"
	"gitee.com/mxcker/go-admin/server/core/internal/modules/sys_menu"
	"gitee.com/mxcker/go-admin/server/core/internal/modules/sys_role"
	"gorm.io/gorm"
)

type serviceStruct struct {
}

func (_ *serviceStruct) GetAdminByUsername(username string) (admin *sys_admin.SysAdmin, err error) {
	admin = &sys_admin.SysAdmin{}
	err = global.DB.Model(sys_admin.SysAdmin{}).Where("username = ?", username).First(admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("管理员不存在")
	}
	return admin, err
}

func (_ *serviceStruct) GetUser(id uint) (user *sys_admin.SysAdmin, err error) {
	user = &sys_admin.SysAdmin{}
	err = global.DB.Preload("Role").Where("id = ?", id).First(user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("管理员账号不存在")
	}
	if user.Status != 1 {
		return nil, errors.New("管理员被禁用")
	}
	var menus []*sys_menu.SysMenu
	if user.Role != nil && user.Role.Default {
		if err = global.DB.Model(&sys_menu.SysMenu{}).Find(&menus).Error; err != nil {
			return nil, err
		}
		user.Role.Menus = menus
	} else {
		return user, global.DB.Model(&sys_role.SysRole{}).Preload("Menus").Find(user.Role).Error
	}
	return user, err
}
