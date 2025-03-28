package sys_menu

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/core/global"
	"gorm.io/gorm/clause"
)

var service = new(Service)

type Service struct {
}

func (s *Service) ListToTree(allList []*SysMenu) (list []*SysMenu) {
	idMap := make(map[uint]*SysMenu)
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

func (s *Service) Get(id uint) (data *SysMenu, err error) {
	data = &SysMenu{}
	err = global.DB.Model(SysMenu{}).Where("id = ?", id).Find(data).Error
	return
}

func (s *Service) List() (list []*SysMenu, total int64, err error) {
	db := global.DB.Model(SysMenu{})
	err = db.Find(&list).Error
	total = int64(len(list))
	return
}

func (s *Service) Create(menu *SysMenu) (err error) {
	err = global.DB.Omit(clause.Associations).Create(menu).Error
	return
}

func (s *Service) Update(data *SysMenu) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return global.DB.Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt", "Children").
		Where("id = ?", data.Id).Updates(data).Error
}

func (s *Service) DeleteByIds(ids []uint) (err error) {
	err = global.DB.Delete(&SysMenu{}, ids).Error
	return
}
