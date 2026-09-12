package sys_role

import (
	"errors"
	"github.com/Gary-Yez/go-admin/internal/permissions"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/Gary-Yez/go-admin/internal/utils"

	request2 "github.com/Gary-Yez/go-admin/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type serviceStruck struct {
}

func (s *serviceStruck) PermissionOptions() ([]*sys_menu.SysMenu, []*ApiOption, error) {
	menus, _, err := sys_menu.Service.List()
	if err != nil {
		return nil, nil, err
	}
	apis := make([]*ApiOption, 0)
	db := state.DB().Table("sys_apis").Clauses(clause.Select{Columns: []clause.Column{
		{Name: "method"}, {Name: "path"}, {Name: "group"}, {Name: "description"},
	}})
	if err := db.Order("id DESC").Scan(&apis).Error; err != nil {
		return nil, nil, err
	}
	return sys_menu.Service.ListToTree(menus), apis, nil
}

func (s *serviceStruck) Get(req *request2.Req) (data *SysRole, err error) {
	data = &SysRole{}
	err = req.WithQuery(state.DB().Model(SysRole{}).Preload("Menus")).First(data).Error
	if err != nil {
		return nil, err
	}
	policy, err := Enforcer.GetFilteredPolicy(0, strconv.Itoa(int(data.Id)))
	if err != nil {
		return nil, err
	}
	pathMaps := make([]*SysCasbinApi, 0)
	for _, v := range policy {
		pathMaps = append(pathMaps, &SysCasbinApi{
			Path:   v[1],
			Method: v[2],
		})
	}
	data.Apis = pathMaps
	return
}

func (s *serviceStruck) List(req *request2.ReqList) (list []*SysRole, total int64, err error) {
	db := req.WithFilter(state.DB().Model(SysRole{}), nil)
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if req.Page == 0 {
		req.Page = 1
	}
	err = req.WithPagination(req.WithSort(db, []string{"id"})).Find(&list).Error
	return
}

func (s *serviceStruck) Create(data *SysRole) (err error) {
	data.IsSuperAdmin = false
	err = state.DB().Omit("Menus.*").Create(data).Error
	return
}

func (s *serviceStruck) Update(data *SysRole) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return state.DB().Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt", "DefaultMenu", "IsSuperAdmin").
		Where("id = ?", data.Id).Updates(data).Error
}

func (s *serviceStruck) DeleteByIds(req *request2.ReqIds) error {
	return state.DB().Transaction(func(tx *gorm.DB) error {
		var roles []SysRole
		if err := req.WithQuery(tx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("is_super_admin = ?", false).Find(&roles).Error; err != nil {
			return err
		}
		ids := make([]uint, 0, len(roles))
		for _, role := range roles {
			ids = append(ids, role.Id)
		}
		if len(ids) == 0 {
			return nil
		}
		var count int64
		if err := tx.Table("sys_admin_role").Where("role_id IN ?", ids).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("该角色正在使用中，请先解除用户关联")
		}
		err := tx.Where("id IN ?", ids).Delete(&SysRole{}).Error
		if err != nil {
			if utils.IsForeignKeyConstraintError(err) {
				return errors.New("该角色正在使用中")
			}
			return err
		}
		for _, v := range ids {
			err = s.DeletePolicy(0, strconv.Itoa(int(v)))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *serviceStruck) UpdatePermission(data *SysRole) (err error) {
	for _, api := range data.Apis {
		if api == nil {
			return errors.New("API权限配置无效")
		}
		if permissions.IsAuthenticatedAPI(api.Method, api.Path) {
			return errors.New("登录后的基础接口由所有角色共享，无需配置")
		}
	}
	roleId := strconv.Itoa(int(data.Id))
	err = s.DeletePolicy(0, roleId)
	if err != nil {
		return err
	}
	var rules [][]string
	//做权限去重处理
	deduplicateMap := make(map[string]bool)
	for _, v := range data.Apis {
		key := roleId + v.Path + v.Method
		if _, ok := deduplicateMap[key]; !ok {
			deduplicateMap[key] = true
			rules = append(rules, []string{roleId, v.Path, v.Method})
		}
	}
	if len(rules) != 0 {
		success, err := Enforcer.AddPolicies(rules)
		if err != nil {
			return err
		}
		if !success {
			return errors.New("API权限修改失败")
		}
	} // 设置空权限无需调用 AddPolicies 方法
	return state.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Menus.*").Model(data).Association("Menus").Replace(data.Menus); err != nil {
			return err
		}
		return tx.Model(data).Update("default_menu", data.DefaultMenu).Error
	})
}

func (s *serviceStruck) DeletePolicy(fieldIndex int, fieldValues ...string) error {
	_, err := Enforcer.RemoveFilteredPolicy(fieldIndex, fieldValues...)
	return err
}
