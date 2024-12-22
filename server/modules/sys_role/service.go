package sys_role

import (
	"errors"
	"fmt"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/models"
	"gorm.io/gorm/clause"
)

type Service struct {
}

func (s *Service) Get(id uint) (data *models.SysRole, err error) {
	data = &models.SysRole{}
	err = global.DB.Model(models.SysRole{}).Where("id = ?", id).Find(data).Error
	return
}

func (s *Service) List() (list []*models.SysRole, total int64, err error) {
	db := global.DB.Model(models.SysRole{}).Preload("Menus")
	err = db.Find(&list).Error
	return
}

func (s *Service) Create(data *models.SysRole) (err error) {
	err = global.DB.Create(data).Error
	return
}

func (s *Service) Update(data *models.SysRole) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return global.DB.Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt").
		Where("id = ?", data.Id).Updates(data).Error
}

func (s *Service) DeleteByIds(ids []uint) (err error) {
	fmt.Println(ids)
	err = global.DB.Where("`default` = ?", false).Delete(&models.SysRole{}, ids).Error
	return
}

func (s *Service) UpdatePermission(data *models.SysRole) (err error) {
	return global.DB.Debug().Model(data).Association("Menus").Replace(data.Menus)
}
