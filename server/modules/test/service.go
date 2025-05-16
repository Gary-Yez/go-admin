package test

import (
    "errors"
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/types/request"
    "gorm.io/gorm/clause"
)

type serviceStruct struct {
}
func (s *serviceStruct) Get(req *request.Req) (data *Test, err error) {
	data = &Test{}
	err = req.WithQuery(global.DB.Model(Test{})).First(data).Error
	return
}

func (s *serviceStruct) List(req *request.ReqList) (list []*Test, total int64, err error) {
	db := req.WithFilter(global.DB.Model(Test{}), nil)
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.WithPagination(req.WithSort(db, nil)).Find(&list).Error
	return
}

func (s *serviceStruct) Create(data *Test) (err error) {
	err = global.DB.
	    Omit(clause.Associations).
	    Create(data).Error
	return
}

func (s *serviceStruct) Update(data *Test) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return global.DB.Select("*").
    		Omit(clause.Associations).
    		Omit("Id", "CreatedAt", "UpdatedAt").
    		Where("id = ?", data.Id).Updates(data).Error
}

func (s *serviceStruct) DeleteByIds(req *request.ReqIds) (err error) {
	return req.WithQuery(global.DB).Delete(&Test{}).Error
}
