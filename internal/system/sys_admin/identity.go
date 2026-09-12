package sys_admin

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"time"

	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_role"
	"gorm.io/gorm"
)

// identity 只保存认证所需的账号状态和仍存在的角色，不包含资料或权限规则。
type identity struct {
	LoginVersion uint64 `json:"login_version"`
	Status       uint   `json:"status"`
	RoleIds      []uint `json:"role_ids"`
}

func identityKey(id uint) string { return "auth:identity:" + strconv.FormatUint(uint64(id), 10) }

func lockIdentity(ctx context.Context, id uint) (cache.Lock, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	lock, err := state.Cache().WaitForLock(ctx, identityKey(id), 30*time.Second, 10*time.Second)
	if err != nil {
		return nil, errors.New("身份缓存繁忙或暂时不可用，请稍后重试")
	}
	return lock, nil
}

func (*serviceStruct) ReadIdentity(id uint) (*identity, error) {
	store := state.Cache()
	key := identityKey(id)
	var data identity
	if err := store.GetJSON(key, &data); err == nil {
		return &data, nil
	} else if !errors.Is(err, cache.ErrCacheNotFound) {
		return nil, errors.New("身份缓存暂时不可用")
	}
	lock, err := lockIdentity(context.Background(), id)
	if err != nil {
		return nil, err
	}
	defer lock.Unlock()
	if err := store.GetJSON(key, &data); err == nil {
		return &data, nil
	} else if !errors.Is(err, cache.ErrCacheNotFound) {
		return nil, errors.New("身份缓存暂时不可用")
	}
	var admin SysAdmin
	if err := state.DB().Select("id", "status", "login_version").First(&admin, id).Error; err != nil {
		return nil, err
	}
	data.Status = admin.Status
	data.LoginVersion = admin.LoginVersion
	if err := state.DB().Model(&SysAdminRole{}).Where("admin_id = ?", id).
		Where("role_id IN (?)", state.DB().Model(&sys_role.SysRole{}).Select("id")).
		Pluck("role_id", &data.RoleIds).Error; err != nil {
		return nil, err
	}
	if err := store.SetJSON(key, data, cache.DefaultTTL); err != nil {
		return nil, errors.New("写入身份缓存失败")
	}
	return &data, nil
}

// 修改与回填持有同一把用户锁，清除缓存后提交事务，避免旧身份被并发回填。
func changeIdentity(ids []uint, change func(*gorm.DB) error) error {
	ids = slices.Clone(ids)
	slices.Sort(ids)
	ids = slices.Compact(ids)
	locks := make([]cache.Lock, 0, len(ids))
	defer func() {
		for i := len(locks) - 1; i >= 0; i-- {
			_ = locks[i].Unlock()
		}
	}()
	for _, id := range ids {
		lock, err := lockIdentity(context.Background(), id)
		if err != nil {
			return err
		}
		locks = append(locks, lock)
	}
	for _, id := range ids {
		if err := state.Cache().Del(identityKey(id)); err != nil {
			return errors.New("清除身份缓存失败，请重试")
		}
	}
	return state.DB().Transaction(change)
}
