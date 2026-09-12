package sys_admin

import (
	"errors"
	"slices"

	"github.com/Gary-Yez/go-admin/internal/system/sys_role"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func setRoleIds(data *SysAdmin) {
	data.RoleIds = make([]uint, 0, len(data.Roles))
	for _, role := range data.Roles {
		data.RoleIds = append(data.RoleIds, role.Id)
	}
}

func prepareRoles(tx *gorm.DB, data *SysAdmin) error {
	ids := data.RoleIds
	unique := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			return errors.New("角色ID无效")
		}
		if !slices.Contains(unique, id) {
			unique = append(unique, id)
		}
	}
	if len(unique) == 0 {
		return errors.New("请至少选择一个角色")
	}
	var roles []*sys_role.SysRole
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ?", unique).Order("id").Find(&roles).Error; err != nil {
		return err
	}
	if len(roles) != len(unique) {
		return errors.New("所选角色不存在，请刷新后重试")
	}
	if !slices.Contains(unique, data.RoleId) {
		return errors.New("请选择已绑定的角色作为默认角色")
	}
	data.RoleIds, data.Roles = unique, roles
	return nil
}

func replaceRoles(tx *gorm.DB, data *SysAdmin) error {
	if err := tx.Where("admin_id = ?", data.Id).Delete(&SysAdminRole{}).Error; err != nil {
		return err
	}
	links := make([]SysAdminRole, 0, len(data.RoleIds))
	for _, id := range data.RoleIds {
		links = append(links, SysAdminRole{AdminId: data.Id, RoleId: id})
	}
	return tx.Create(&links).Error
}
