package sys_admin

import (
	"errors"

	"github.com/Gary-Yez/go-admin/internal/system/sys_role"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 必须在读取或修改账号之前调用；角色行锁持有到事务结束，串行化多实例的管理员变更。
func lockSuperAdminRoles(tx *gorm.DB) ([]uint, error) {
	var roles []sys_role.SysRole
	if err := tx.Select("id").Where("is_super_admin = ?", true).Order("id").
		Clauses(clause.Locking{Strength: "UPDATE"}).Find(&roles).Error; err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, errors.New("超级管理员角色不存在，请检查系统初始化数据")
	}
	ids := make([]uint, 0, len(roles))
	for _, role := range roles {
		ids = append(ids, role.Id)
	}
	return ids, nil
}

// 在同一事务内检查修改后的数据；返回错误时由调用方回滚账号和角色关联变更。
func ensureSuperAdmin(tx *gorm.DB, roleIds []uint) error {
	var admin SysAdmin
	roles := tx.Model(&SysAdminRole{}).Select("admin_id").Where("role_id IN ?", roleIds)
	err := tx.Select("id").Where("status = ?", 1).Where("id IN (?)", roles).
		Clauses(clause.Locking{Strength: "UPDATE"}).Take(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("系统至少需要保留一个启用且拥有超级管理员角色的账号")
	}
	return err
}
