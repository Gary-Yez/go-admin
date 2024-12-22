package sys_admin

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/models"
	"gitee.com/mxcker/go-admin/server/models/request"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

type Service struct {
}

func (s *Service) GeneratePassHash(data *models.SysAdmin) error {
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	data.PasswordHash = string(password)
	return nil
}

func (s *Service) Get(id uint) (data *models.SysAdmin, err error) {
	data = &models.SysAdmin{}
	err = global.DB.Model(models.SysAdmin{}).Where("id = ?", id).Find(data).Error
	return
}

func (s *Service) List(req *request.ReqList) (list []*models.SysAdmin, total int64, err error) {
	db := global.DB.Model(models.SysAdmin{})
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	req.SetDB(db)
	err = db.Find(&list).Error
	return
}

func (s *Service) Create(data *models.SysAdmin) (err error) {
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

func (s *Service) Update(data *models.SysAdmin) (err error) {
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

func (s *Service) DeleteByIds(ids []uint) (err error) {
	err = global.DB.Debug().Model(&models.SysAdmin{}).Where("`default` = ?", false).Delete(&models.SysAdmin{}, ids).Error
	return
}
