package sys_role

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gary-Yez/go-admin/dberror"
	"github.com/Gary-Yez/go-admin/internal/permissions"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/Gary-Yez/go-admin/internal/utils"
	gormadapter "github.com/casbin/gorm-adapter/v3"

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
	return dberror.Unique(err, data)
}

func (s *serviceStruck) Update(data *SysRole) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	err = state.DB().Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt", "DefaultMenu", "IsSuperAdmin").
		Where("id = ?", data.Id).Updates(data).Error
	return dberror.Unique(err, data)
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

func (s *serviceStruck) UpdatePermission(data *SysRole) error {
	if data == nil || data.Id == 0 {
		return errors.New("请选择角色")
	}
	roleId := strconv.Itoa(int(data.Id))
	var rules [][]string
	seen := make(map[[2]string]bool)
	for _, api := range data.Apis {
		if api == nil || api.Path == "" || api.Method == "" {
			return errors.New("API权限配置无效")
		}
		if permissions.IsAuthenticatedAPI(api.Method, api.Path) {
			return errors.New("登录后的基础接口由所有角色共享，无需配置")
		}
		key := [2]string{api.Path, api.Method}
		if !seen[key] {
			seen[key] = true
			rules = append(rules, []string{roleId, api.Path, api.Method})
		}
	}
	adapter, ok := Enforcer.GetAdapter().(*gormadapter.Adapter)
	if !ok {
		return errors.New("角色权限存储未就绪")
	}
	transaction, err := adapter.BeginTransaction(context.Background())
	if err != nil {
		return err
	}
	defer transaction.Rollback()
	txAdapter := transaction.GetAdapter().(*gormadapter.Adapter)
	// 策略、菜单关联和默认首页共用事务；角色行锁串行化同一角色的权限修改。
	tx := txAdapter.GetDb().Session(&gorm.Session{NewDB: true})
	var role SysRole
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&role, data.Id).Error; err != nil {
		return err
	}
	if err := txAdapter.RemoveFilteredPolicy("p", "p", 0, roleId); err != nil {
		return err
	}
	if len(rules) > 0 {
		if err := txAdapter.AddPolicies("p", "p", rules); err != nil {
			return err
		}
	}
	if err := tx.Omit("Menus.*").Model(&role).Association("Menus").Replace(data.Menus); err != nil {
		return err
	}
	if err := tx.Model(&role).Update("default_menu", data.DefaultMenu).Error; err != nil {
		return err
	}
	if err := transaction.Commit(); err != nil {
		return err
	}
	// 提交后再刷新本机缓存并通知其他实例，避免重复写库或提前公布未提交权限。
	if err := s.ReloadPolicy(); err != nil {
		return fmt.Errorf("角色权限已保存，但权限缓存更新失败：%w", err)
	}
	return nil
}

func (s *serviceStruck) DeletePolicy(fieldIndex int, fieldValues ...string) error {
	_, err := Enforcer.RemoveFilteredPolicy(fieldIndex, fieldValues...)
	return err
}
