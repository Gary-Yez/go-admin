package sys_global_variable

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/core/global"
	request2 "gitee.com/mxcker/go-admin/server/core/types/request"
	"gorm.io/gorm/clause"
)

type serviceStruct struct {
}

func (s *serviceStruct) Get(req *request2.Req) (data *SysGlobalVariable, err error) {
	data = &SysGlobalVariable{}
	err = req.BuildQuery(global.DB.Model(SysGlobalVariable{})).First(data).Error
	return
}

func (s *serviceStruct) List(req *request2.ReqList) (list []*SysGlobalVariable, total int64, err error) {
	db := global.DB.Model(SysGlobalVariable{})
	err = req.BuildWhere(db).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.BuildQuery(db).Find(&list).Error
	return
}

func (s *serviceStruct) Create(data *SysGlobalVariable) (err error) {
	err = global.DB.
		Omit(clause.Associations).
		Create(data).Error
	return
}

func (s *serviceStruct) Update(data *SysGlobalVariable) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return global.DB.Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt").
		Where("id = ?", data.Id).Updates(data).Error
}

func (s *serviceStruct) DeleteByIds(req *request2.ReqIds) (err error) {
	return req.BuildQuery(global.DB).Delete(&SysGlobalVariable{}).Error
}
