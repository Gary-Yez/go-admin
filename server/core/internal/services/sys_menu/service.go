package sys_menu

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/utils"
	"gitee.com/mxcker/go-admin/server/utils/request"
	"gorm.io/gorm/clause"
)

type serviceStruct struct {
}

func (s *serviceStruct) ListToTree(allList []*SysMenu) (list []*SysMenu) {
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

func (s *serviceStruct) Get(req *request.Req) (data *SysMenu, err error) {
	data = &SysMenu{}
	err = req.WithQuery(global.DB.Model(SysMenu{})).First(data).Error
	return
}

func (s *serviceStruct) List() (list []*SysMenu, total int64, err error) {
	db := global.DB.Model(SysMenu{})
	err = db.Find(&list).Error
	total = int64(len(list))
	return
}

func (s *serviceStruct) Create(menu *SysMenu) (err error) {
	err = global.DB.Omit(clause.Associations).Create(menu).Error
	return
}

func (s *serviceStruct) Update(data *SysMenu) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return global.DB.Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt", "Children").
		Where("id = ?", data.Id).Updates(data).Error
}

func (s *serviceStruct) DeleteByIds(req *request.ReqIds) (err error) {
	err = req.WithQuery(global.DB).Delete(&SysMenu{}).Error
	if utils.IsForeignKeyConstraintError(err) {
		return errors.New("该菜单存在子菜单，请先删除子菜单")
	}
	return
}
