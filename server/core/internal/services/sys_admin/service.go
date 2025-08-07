package sys_admin

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/utils/request"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

type serviceStruct struct {
}

func (s *serviceStruct) GeneratePassHash(data *SysAdmin) error {
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	data.PasswordHash = string(password)
	return nil
}

func (s *serviceStruct) Get(req *request.Req) (data *SysAdmin, err error) {
	data = &SysAdmin{}
	err = req.WithQuery(global.DB.Model(SysAdmin{})).First(data).Error
	return
}

func (s *serviceStruct) List(req *request.ReqList) (list []*SysAdmin, total int64, err error) {
	db := req.WithFilter(global.DB.Model(SysAdmin{}).Omit("ApiToken"), nil)
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.WithPagination(req.WithSort(db, nil)).Find(&list).Error
	return
}

func (s *serviceStruct) Create(data *SysAdmin) (err error) {
	if data.Password == "" {
		data.Password = "123456"
	}
	err = s.GeneratePassHash(data)
	if err != nil {
		return err
	}
	err = global.DB.Omit(clause.Associations).Create(data).Error
	return
}

func (s *serviceStruct) Update(data *SysAdmin) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	query := global.DB.Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt", "ApiToken")
	if data.Password != "" {
		err = s.GeneratePassHash(data)
		if err != nil {
			return err
		}
	} else {
		query = query.Omit("PasswordHash").Updates(data)
	}
	return query.
		Where("id = ?", data.Id).Updates(data).Error
}

func (s *serviceStruct) DeleteByIds(req *request.ReqIds) (err error) {
	err = req.WithQuery(global.DB.Model(&SysAdmin{}).Debug()).Where("`default` = ?", false).Delete(&SysAdmin{}).Error
	return
}
