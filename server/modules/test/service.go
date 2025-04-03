package test

import (
    "errors"
    "gorm.io/gorm/clause"
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/request"
	
)

type serviceStruct struct {
}
func (s *serviceStruct) Get(id uint) (data *Test, err error) {
	data = &Test{}
	err = global.DB.Model(Test{}).Where("id = ?", id).Find(data).Error
	return
}

func (s *serviceStruct) List(req *request.ReqList) (list []*Test, total int64, err error) {
	db := global.DB.Model(Test{})
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	req.SetDB(db)
	err = db.Find(&list).Error
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

func (s *serviceStruct) DeleteByIds(ids []uint) (err error) {
	return global.DB.Delete(&Test{}, ids).Error
}
