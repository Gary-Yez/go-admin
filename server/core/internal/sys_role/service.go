package sys_role

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/core/global"
	request2 "gitee.com/mxcker/go-admin/server/core/types/request"
	"gorm.io/gorm/clause"
)

type serviceStruck struct {
}

func (s *serviceStruck) Get(req *request2.Req) (data *SysRole, err error) {
	data = &SysRole{}
	err = req.BuildQuery(global.DB.Model(SysRole{})).First(data).Error
	return
}

func (s *serviceStruck) List() (list []*SysRole, total int64, err error) {
	db := global.DB.Model(SysRole{}).Preload("Menus")
	err = db.Find(&list).Error
	return
}

func (s *serviceStruck) Create(data *SysRole) (err error) {
	err = global.DB.Create(data).Error
	return
}

func (s *serviceStruck) Update(data *SysRole) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return global.DB.Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt").
		Where("id = ?", data.Id).Updates(data).Error
}

func (s *serviceStruck) DeleteByIds(req *request2.ReqIds) (err error) {
	err = req.BuildQuery(global.DB).Where("`default` = ?", false).Delete(&SysRole{}).Error
	return
}

func (s *serviceStruck) UpdatePermission(data *SysRole) (err error) {
	return global.DB.Debug().Model(data).Association("Menus").Replace(data.Menus)
}
