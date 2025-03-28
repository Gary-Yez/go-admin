package sys_admin

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/request"
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

func (s *serviceStruct) Get(id uint) (data *SysAdmin, err error) {
	data = &SysAdmin{}
	err = global.DB.Model(SysAdmin{}).Where("id = ?", id).Find(data).Error
	return
}

func (s *serviceStruct) List(req *request.ReqList) (list []*SysAdmin, total int64, err error) {
	db := global.DB.Model(SysAdmin{})
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	req.SetDB(db)
	err = db.Find(&list).Error
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
		Omit("Id", "CreatedAt", "UpdatedAt")
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

func (s *serviceStruct) DeleteByIds(ids []uint) (err error) {
	err = global.DB.Debug().Model(&SysAdmin{}).Where("`default` = ?", false).Delete(&SysAdmin{}, ids).Error
	return
}
