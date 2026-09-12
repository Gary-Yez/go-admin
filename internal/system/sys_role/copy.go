package sys_role

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Gary-Yez/go-admin/internal/permissions"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *serviceStruck) Copy(body *CopyBody) (*SysRole, error) {
	name := strings.TrimSpace(body.Name)
	if body.Id == 0 || name == "" {
		return nil, errors.New("请选择来源角色并填写新角色名称")
	}
	adapter, ok := Enforcer.GetAdapter().(*gormadapter.Adapter)
	if !ok {
		return nil, errors.New("角色权限存储未就绪")
	}
	transaction, err := adapter.BeginTransaction(context.Background())
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()
	txAdapter := transaction.GetAdapter().(*gormadapter.Adapter)
	// 清除适配器绑定的权限表名，保留同一个数据库事务连接。
	tx := txAdapter.GetDb().Session(&gorm.Session{NewDB: true})
	var source SysRole
	if err := tx.Preload("Menus").First(&source, body.Id).Error; err != nil {
		return nil, err
	}
	var count int64
	if err := tx.Model(&SysRole{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("角色名称已存在")
	}
	role := &SysRole{Name: name, Menus: source.Menus, DefaultMenu: source.DefaultMenu}
	var sourceRules [][]string
	if source.IsSuperAdmin {
		// 超级管理员副本是当前权限快照，不复制通配权限或系统角色标识。
		if err := tx.Find(&role.Menus).Error; err != nil {
			return nil, err
		}
		for _, route := range state.AdminRoutes() {
			path := strings.TrimPrefix(route.Path, state.Config().Server.ApiPrefix)
			if permissions.IsAuthenticatedAPI(route.Method, path) {
				continue
			}
			sourceRules = append(sourceRules, []string{"", path, route.Method})
		}
	} else {
		sourceRules, err = Enforcer.GetFilteredPolicy(0, strconv.Itoa(int(source.Id)))
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Omit(clause.Associations).Create(role).Error; err != nil {
		return nil, err
	}
	// 只写关联，菜单记录本身保持独立管理。
	menus := role.Menus
	role.Menus = []*sys_menu.SysMenu{}
	if err := tx.Omit("Menus.*").Model(role).Association("Menus").Replace(menus); err != nil {
		return nil, err
	}
	rules := make([][]string, 0, len(sourceRules))
	for _, rule := range sourceRules {
		if len(rule) < 3 {
			continue
		}
		rules = append(rules, []string{strconv.Itoa(int(role.Id)), rule[1], rule[2]})
	}
	if len(rules) > 0 {
		if err := txAdapter.AddPolicies("p", "p", rules); err != nil {
			return nil, err
		}
	}
	if err := transaction.Commit(); err != nil {
		return nil, err
	}
	// 数据提交后更新本机权限及 Casbin watcher，避免其他实例继续使用旧缓存。
	if len(rules) > 0 {
		if err := s.ReloadPolicy(); err != nil {
			return role, fmt.Errorf("角色已拷贝，但权限缓存更新失败：%w", err)
		}
	}
	return role, nil
}
