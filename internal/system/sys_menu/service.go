package sys_menu

import (
	"errors"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/utils"
	"slices"
	"strings"

	request2 "github.com/Gary-Yez/go-admin/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type serviceStruct struct {
}

// 为已授权菜单补齐路由目录，只加入祖先目录，不加入同级菜单或未授权页面。
func (s *serviceStruct) WithAncestors(menus []*SysMenu) ([]*SysMenu, error) {
	result := append([]*SysMenu{}, menus...)
	seen := make(map[uint]bool, len(menus))
	for _, menu := range menus {
		seen[menu.Id] = true
	}
	for pending := menus; len(pending) > 0; {
		ids := make([]uint, 0)
		for _, menu := range pending {
			if menu.ParentId != nil && !seen[*menu.ParentId] {
				seen[*menu.ParentId] = true
				ids = append(ids, *menu.ParentId)
			}
		}
		if len(ids) == 0 {
			break
		}
		var parents []*SysMenu
		if err := state.DB().Where("id IN ? AND (component = ? OR component IS NULL)", ids, "").Find(&parents).Error; err != nil {
			return nil, err
		}
		result = append(result, parents...)
		pending = parents
	}
	return result, nil
}

func (s *serviceStruct) ListToTree(allList []*SysMenu) (list []*SysMenu) {
	list = make([]*SysMenu, 0)
	idMap := make(map[uint]*SysMenu)
	for _, v := range allList {
		v.Children = nil
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

func (s *serviceStruct) List() (list []*SysMenu, total int64, err error) {
	db := state.DB().Model(SysMenu{})
	err = db.Order("sort ASC").Order("id ASC").Find(&list).Error
	total = int64(len(list))
	return
}

// Sort 在事务中移动整棵子树，并重新排列源父级和目标父级下的菜单。
func (s *serviceStruct) Sort(req *SortBody) error {
	if req == nil || req.Id == 0 || req.TargetId == 0 || req.Id == req.TargetId || !slices.Contains([]string{"before", "after", "inside"}, req.Position) {
		return errors.New("拖拽位置无效")
	}
	return state.DB().Transaction(func(tx *gorm.DB) error {
		var menus []SysMenu
		// 统一按 ID 加锁，防止两个并发跨层级移动相互构成循环。
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Order("id ASC").Find(&menus).Error; err != nil {
			return err
		}
		byID := make(map[uint]*SysMenu, len(menus))
		for i := range menus {
			byID[menus[i].Id] = &menus[i]
		}
		moved, target := byID[req.Id], byID[req.TargetId]
		if moved == nil || target == nil {
			return errors.New("菜单已不存在，请刷新后重试")
		}
		parentID := func(id *uint) uint {
			if id == nil {
				return 0
			}
			return *id
		}
		newParent := target.ParentId
		if req.Position == "inside" {
			newParent = &target.Id
		}
		oldParentID, newParentID := parentID(moved.ParentId), parentID(newParent)
		if oldParentID != parentID(req.SourceParentId) || newParentID != parentID(req.TargetParentId) {
			return errors.New("菜单层级已变化，请刷新后重试")
		}
		if newParentID != 0 {
			parent := byID[newParentID]
			if parent == nil || parent.Component != "" {
				return errors.New("只能移入路由目录，页面组件不能包含子菜单")
			}
		}
		seen := map[uint]bool{}
		for id := newParentID; id != 0; {
			if id == moved.Id || seen[id] {
				return errors.New("不能移入自身或后代菜单")
			}
			seen[id] = true
			parent := byID[id]
			if parent == nil {
				return errors.New("父菜单不存在，请刷新后重试")
			}
			id = parentID(parent.ParentId)
		}
		siblings := func(parent uint) []uint {
			ids := []uint{}
			for _, menu := range menus {
				if parentID(menu.ParentId) == parent {
					ids = append(ids, menu.Id)
				}
			}
			slices.SortFunc(ids, func(a, b uint) int {
				if byID[a].Sort < byID[b].Sort {
					return -1
				}
				if byID[a].Sort > byID[b].Sort {
					return 1
				}
				if a < b {
					return -1
				}
				if a > b {
					return 1
				}
				return 0
			})
			return ids
		}
		source, destination := siblings(oldParentID), siblings(newParentID)
		if !slices.Equal(source, req.SourceIds) || !slices.Equal(destination, req.TargetIds) {
			return errors.New("菜单顺序已变化，请刷新后重试")
		}
		if oldParentID != newParentID {
			for _, id := range destination {
				if strings.Trim(byID[id].Path, "/") == strings.Trim(moved.Path, "/") {
					return errors.New("目标父级下已存在相同路由地址的菜单")
				}
			}
		}
		source = slices.DeleteFunc(source, func(id uint) bool { return id == moved.Id })
		destination = slices.DeleteFunc(destination, func(id uint) bool { return id == moved.Id })
		index := len(destination)
		if req.Position != "inside" {
			index = slices.Index(destination, target.Id)
			if index < 0 {
				return errors.New("目标菜单已变化，请刷新后重试")
			}
			if req.Position == "after" {
				index++
			}
		}
		destination = slices.Insert(destination, index, moved.Id)
		if err := tx.Model(&SysMenu{}).Where("id = ?", moved.Id).Update("parent_id", newParent).Error; err != nil {
			return err
		}
		groups := [][]uint{destination}
		if oldParentID != newParentID {
			groups = append(groups, source)
		}
		for _, ids := range groups {
			for sort, id := range ids {
				if err := tx.Model(&SysMenu{}).Where("id = ?", id).Update("sort", sort).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *serviceStruct) Create(menu *SysMenu) (err error) {
	err = state.DB().Omit(clause.Associations).Create(menu).Error
	return
}

func (s *serviceStruct) Update(data *SysMenu) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return state.DB().Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt", "Children").
		Where("id = ?", data.Id).Updates(data).Error
}

func (s *serviceStruct) DeleteByIds(req *request2.ReqIds) (err error) {
	err = req.WithQuery(state.DB()).Delete(&SysMenu{}).Error
	if utils.IsForeignKeyConstraintError(err) {
		return errors.New("该菜单存在子菜单，请先删除子菜单")
	}
	return
}
