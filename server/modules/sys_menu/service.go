package sys_menu

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/models"
	"gorm.io/gorm/clause"
)

type Service struct {
}

func (s *Service) ListToTree(allList []*models.SysMenu) (list []*models.SysMenu) {
	idMap := make(map[uint]*models.SysMenu)
	for _, v := range allList {
		idMap[v.Id] = v
	}
	for _, v := range allList {
		if v.ParentId != nil {
			parentNode, exists := idMap[*v.ParentId]
			if exists {
				parentNode.Children = append(parentNode.Children, idMap[v.Id])
			}
		} else {
			list = append(list, idMap[v.Id])
		}
	}
	return list
}

func (s *Service) Get(id uint) (data *models.SysMenu, err error) {
	data = &models.SysMenu{}
	err = global.DB.Model(models.SysMenu{}).Where("id = ?", id).Find(data).Error
	return
}

func (s *Service) List() (list []*models.SysMenu, total int64, err error) {
	db := global.DB.Model(models.SysMenu{})
	err = db.Find(&list).Error
	total = int64(len(list))
	return
}

func (s *Service) Create(menu *models.SysMenu) (err error) {
	err = global.DB.Omit(clause.Associations).Create(menu).Error
	return
}

func (s *Service) Update(data *models.SysMenu) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return global.DB.Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt", "Children").
		Where("id = ?", data.Id).Updates(data).Error
}

func (s *Service) DeleteByIds(ids []uint) (err error) {
	err = global.DB.Delete(&models.SysMenu{}, ids).Error
	return
}
