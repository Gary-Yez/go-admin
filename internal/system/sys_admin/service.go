package sys_admin

import (
	"context"
	"errors"
	"github.com/Gary-Yez/go-admin/dberror"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_role"
	"github.com/Gary-Yez/go-admin/internal/utils"
	"slices"

	request2 "github.com/Gary-Yez/go-admin/request"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type serviceStruct struct {
}

func (s *serviceStruct) RoleOptions() ([]RoleOption, error) {
	options := make([]RoleOption, 0)
	err := state.DB().Model(&sys_role.SysRole{}).Select("id", "name").Order("id").Scan(&options).Error
	return options, err
}

func (s *serviceStruct) GeneratePassHash(data *SysAdmin) error {
	if err := utils.ValidatePassword(data.Password); err != nil {
		return err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	data.PasswordHash = string(password)
	return nil
}

func (s *serviceStruct) List(req *AdminListQuery) (list []*SysAdmin, total int64, err error) {
	db := req.WithFilter(state.DB().Model(SysAdmin{}), []string{"username", "nickname", "phone", "email", "status"})
	if req.RoleId != 0 {
		members := state.DB().Model(&SysAdminRole{}).Select("admin_id").Where("role_id = ?", req.RoleId)
		db = db.Where("id IN (?)", members)
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.WithPagination(req.WithSort(db.Preload("Roles"), []string{"id"})).Find(&list).Error
	for _, admin := range list {
		setRoleIds(admin)
	}
	return
}

func (s *serviceStruct) Create(data *SysAdmin) (err error) {
	err = s.GeneratePassHash(data)
	if err != nil {
		return err
	}
	return state.DB().Transaction(func(tx *gorm.DB) error {
		if err := prepareRoles(tx, data); err != nil {
			return err
		}
		if err := tx.Omit(clause.Associations).Create(data).Error; err != nil {
			return dberror.Unique(err, data)
		}
		return replaceRoles(tx, data)
	})
}

func (s *serviceStruct) Update(data *SysAdmin) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	if data.Password != "" {
		err = s.GeneratePassHash(data)
		if err != nil {
			return err
		}
	}
	// 先取得身份锁，再比较数据库中的实际值，避免并发编辑绕过缓存失效。
	lock, err := lockIdentity(context.Background(), data.Id)
	if err != nil {
		return err
	}
	defer lock.Unlock()
	return state.DB().Transaction(func(tx *gorm.DB) error {
		superRoleIds, err := lockSuperAdminRoles(tx)
		if err != nil {
			return err
		}
		var existing SysAdmin
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&existing, data.Id).Error; err != nil {
			return err
		}
		var roleIds []uint
		if err := tx.Model(&SysAdminRole{}).Where("admin_id = ?", data.Id).Order("role_id").Pluck("role_id", &roleIds).Error; err != nil {
			return err
		}
		data.RoleIds = slices.Clone(data.RoleIds)
		slices.Sort(data.RoleIds)
		data.RoleIds = slices.Compact(data.RoleIds)
		rolesChanged := !slices.Equal(roleIds, data.RoleIds)
		if rolesChanged {
			if err := prepareRoles(tx, data); err != nil {
				return err
			}
		} else if !slices.Contains(roleIds, data.RoleId) {
			return errors.New("请选择已绑定的角色作为默认角色")
		}
		if existing.Status != data.Status || rolesChanged || data.Password != "" {
			if err := state.Cache().Del(identityKey(data.Id)); err != nil {
				return errors.New("清除身份缓存失败，请重试")
			}
		}
		data.LoginVersion = existing.LoginVersion
		if data.Password != "" {
			data.LoginVersion++
		}
		query := tx.Model(&SysAdmin{}).Select("*").Omit(clause.Associations).Omit("Id", "CreatedAt", "UpdatedAt")
		if data.Password == "" {
			query = query.Omit("PasswordHash")
		}
		if err := query.Where("id = ?", data.Id).Updates(data).Error; err != nil {
			return dberror.Unique(err, data)
		}
		if rolesChanged {
			if err := replaceRoles(tx, data); err != nil {
				return err
			}
		}
		if existing.Status != data.Status || rolesChanged {
			return ensureSuperAdmin(tx, superRoleIds)
		}
		return nil
	})
}

func (s *serviceStruct) DeleteByIds(req *request2.ReqIds) (err error) {
	if err := req.Validate(); err != nil {
		return err
	}
	return changeIdentity(req.Ids, func(tx *gorm.DB) error {
		superRoleIds, err := lockSuperAdminRoles(tx)
		if err != nil {
			return err
		}
		if err := tx.Where("admin_id IN ?", req.Ids).Delete(&SysAdminRole{}).Error; err != nil {
			return err
		}
		if err := req.WithQuery(tx).Delete(&SysAdmin{}).Error; err != nil {
			return err
		}
		return ensureSuperAdmin(tx, superRoleIds)
	})
}

// 密码和登录版本原子更新，并与身份缓存回填共用锁。
func (s *serviceStruct) ChangePassword(id uint, oldHash, newHash string, version uint64) error {
	return changeIdentity([]uint{id}, func(tx *gorm.DB) error {
		result := tx.Model(&SysAdmin{}).Where("id = ? AND password_hash = ? AND login_version = ?", id, oldHash, version).
			Updates(map[string]any{"password_hash": newHash, "login_version": version + 1})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("登录或密码已变更，请重新登录")
		}
		return nil
	})
}
